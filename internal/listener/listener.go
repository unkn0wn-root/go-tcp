package listener

import (
	"fmt"
	"net"
	"sync"
	"syscall"
)

// CustomListener implements a basic TCP listener
type CustomListener struct {
	socket     int
	sockaddr   syscall.Sockaddr
	closed     bool
	closeMutex sync.RWMutex
	acceptChan chan net.Conn
}

// NewCustomListener creates a new TCP listener
// @toDo: Add support for custom listener options
func NewCustomListener(address string, port int) (*CustomListener, error) {
	// Create TCP socket
	socket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to create socket: %v", err)
	}

	// Set socket options
	err = syscall.SetsockoptInt(socket, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	if err != nil {
		syscall.Close(socket)
		return nil, fmt.Errorf("failed to set socket options: %v", err)
	}

	// Convert IP address
	ip := net.ParseIP(address)
	if ip == nil {
		syscall.Close(socket)
		return nil, fmt.Errorf("invalid IP address")
	}

	// Create sockaddr_in structure
	sockaddr := &syscall.SockaddrInet4{Port: port}
	copy(sockaddr.Addr[:], ip.To4())

	// Bind socket
	if err := syscall.Bind(socket, sockaddr); err != nil {
		syscall.Close(socket)
		return nil, fmt.Errorf("failed to bind: %v", err)
	}

	// Listen
	if err := syscall.Listen(socket, syscall.SOMAXCONN); err != nil {
		syscall.Close(socket)
		return nil, fmt.Errorf("failed to listen: %v", err)
	}

	listener := &CustomListener{
		socket:     socket,
		sockaddr:   sockaddr,
		acceptChan: make(chan net.Conn),
	}

	go listener.acceptLoop()

	return listener, nil
}

// Accept implements the net.Listener Accept method
func (l *CustomListener) Accept() (net.Conn, error) {
	conn := <-l.acceptChan
	return conn, nil
}

// Close implements the net.Listener Close method
func (l *CustomListener) Close() error {
	l.closeMutex.Lock()
	defer l.closeMutex.Unlock()

	if l.closed {
		return nil
	}

	l.closed = true
	return syscall.Close(l.socket)
}

// Addr implements the net.Listener Addr method
func (l *CustomListener) Addr() net.Addr {
	sa4, ok := l.sockaddr.(*syscall.SockaddrInet4)
	if !ok {
		return nil
	}

	ip := net.IPv4(sa4.Addr[0], sa4.Addr[1], sa4.Addr[2], sa4.Addr[3])
	return &net.TCPAddr{IP: ip, Port: sa4.Port}
}

func (l *CustomListener) acceptLoop() {
	for {
		l.closeMutex.RLock()
		if l.closed {
			l.closeMutex.RUnlock()
			return
		}
		l.closeMutex.RUnlock()

		nfd, sa, err := syscall.Accept(l.socket)
		if err != nil {
			if l.closed {
				return
			}
			continue
		}

		var remoteAddr net.Addr
		if sa4, ok := sa.(*syscall.SockaddrInet4); ok {
			ip := net.IPv4(sa4.Addr[0], sa4.Addr[1], sa4.Addr[2], sa4.Addr[3])
			remoteAddr = &net.TCPAddr{IP: ip, Port: sa4.Port}
		}

		conn := NewCustomConn(nfd, remoteAddr, l.Addr())
		l.acceptChan <- conn
	}
}
