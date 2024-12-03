package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/unkn0wn-root/go-tcp/pkg/logger"
)

type TCPServer struct {
	listener     net.Listener
	address      string
	port         int
	connections  sync.Map
	logger       *logger.Logger
	readTimeout  time.Duration
	writeTimeout time.Duration
	running      atomic.Bool
}

func NewTCPServer(address string, port int, readTimeout, writeTimeout time.Duration) *TCPServer {
	srv := &TCPServer{
		address:      address,
		port:         port,
		logger:       logger.NewLogger(),
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
	}

	srv.running.Store(true)
	return srv
}

// Start begins listening for incoming connections
func (s *TCPServer) Start() error {
	addr := fmt.Sprintf("%s:%d", s.address, s.port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}

	s.listener = listener
	s.logger.Info("Server listening on %s", addr)

	return s.acceptConnections()
}

func (s *TCPServer) Stop() error {
	s.running.Store(false)
	if s.listener != nil {
		s.logger.Info("Shutting down server...")
		s.connections.Range(func(key, value interface{}) bool {
			if conn, ok := value.(net.Conn); ok {
				conn.Close()
			}
			return true
		})
		return s.listener.Close()
	}
	return nil
}

func (s *TCPServer) acceptConnections() error {
	for s.running.Load() {
		conn, err := s.listener.Accept()
		if err != nil {
			if !s.running.Load() {
				return nil
			}

			if _, ok := err.(net.Error); ok {
				s.logger.Error("Temporary error accepting connection: %v", err)
				continue
			}
			return err
		}

		s.connections.Store(conn.RemoteAddr().String(), conn)
		go s.handleConnection(conn)
	}
	return nil
}

func (s *TCPServer) handleConnection(conn net.Conn) {
	defer func() {
		if r := recover(); r != nil {
			s.logger.Error("Recovered from panic in connection handler: %v", r)
		}
		conn.Close()
		s.connections.Delete(conn.RemoteAddr().String())
	}()

	remoteAddr := conn.RemoteAddr().String()
	s.logger.Info("New connection from %s", remoteAddr)

	reader := bufio.NewReader(conn)

	for {
		// Set read deadline
		if err := conn.SetReadDeadline(time.Now().Add(s.readTimeout)); err != nil {
			s.logger.Error("Failed to set read deadline: %v", err)
			return
		}

		message, err := reader.ReadString('\n')
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				s.logger.Debug("Read timeout for %s", remoteAddr)
				continue
			}
			s.logger.Error("Client %s disconnected: %v", remoteAddr, err)
			return
		}

		message = strings.TrimSpace(message)
		s.logger.Info("Received from %s: %s", remoteAddr, message)

		// Set write deadline
		if err := conn.SetWriteDeadline(time.Now().Add(s.writeTimeout)); err != nil {
			s.logger.Error("Failed to set write deadline: %v", err)
			return
		}

		response := fmt.Sprintf("Server received: %s\n", message)
		_, err = conn.Write([]byte(response))
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				s.logger.Error("Write timeout for %s", remoteAddr)
				return
			}
			s.logger.Error("Failed to send response to %s: %v", remoteAddr, err)
			return
		}

		if strings.ToLower(message) == "quit" {
			s.logger.Info("Client %s requested to quit", remoteAddr)
			return
		}
	}
}
