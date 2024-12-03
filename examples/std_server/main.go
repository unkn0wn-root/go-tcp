package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/unkn0wn-root/go-tcp/internal/server"
)

func main() {
	read, write := 30*time.Second, 30*time.Second
	server := server.NewTCPServer("localhost", 8080, read, write)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	log.Println("Server started. Press Ctrl+C to stop")

	<-sigChan
	log.Println("Shutting down server...")

	if err := server.Stop(); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}
}
