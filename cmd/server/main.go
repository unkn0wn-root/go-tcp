package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/unkn0wn-root/go-tcp/internal/server"
)

var (
	serverType   string
	address      string
	port         int
	readTimeout  time.Duration
	writeTimeout time.Duration
)

func init() {
	flag.StringVar(&serverType, "type", "standard", "server type (standard or raw)")
	flag.StringVar(&address, "addr", "localhost", "server address")
	flag.IntVar(&port, "port", 8080, "server port")
	flag.DurationVar(&readTimeout, "read-timeout", 30*time.Second, "read timeout duration (e.g., 30s)")
	flag.DurationVar(&writeTimeout, "write-timeout", 30*time.Second, "write timeout duration (e.g., 30s)")
}

func main() {
	flag.Parse()

	var srv interface {
		Start() error
		Stop() error
	}

	switch serverType {
	case "standard":
		srv = server.NewTCPServer(address, port, readTimeout, writeTimeout)
	case "raw":
		srv = server.NewRawTCPServer(address, port, readTimeout, writeTimeout)
	default:
		log.Fatalf("Unknown server type: %s", serverType)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	log.Printf("Server started on %s:%d. Press Ctrl+C to stop", address, port)

	<-sigChan
	log.Println("Shutting down server...")

	if err := srv.Stop(); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}
}
