package server

import (
	"fmt"
	"log"
	"net"
)

// StartTCPServer starts a TCP server that listens on the specified port.
func StartTCPServer(port int) {
	// Create the address string for the server: "0.0.0.0:port"
	address := fmt.Sprintf("0.0.0.0:%d", port)

	// Listen on the specified port using TCP
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Error starting TCP server on port %d: %v", port, err)
	}
	defer listener.Close()

	fmt.Printf("TCP server started on %s...\n", address)

	// Accept incoming connections in a loop
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		// Handle each connection in a separate goroutine for concurrency
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	fmt.Printf("Client connected from %s\n", conn.RemoteAddr().String())

	// Handle SMTP commands for this session
	HandleSMTPCommands(conn)
}
