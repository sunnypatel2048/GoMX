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

// handleConnection handles the interaction with a single client.
func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Print the remote address of the client.
	fmt.Printf("Client connected from %s\n", conn.RemoteAddr().String())

	// Send a basic SMTP greeting to the client.
	greeting := "220 Welcome to GoMX Mail Server\r\n"
	_, err := conn.Write([]byte(greeting))
	if err != nil {
		log.Printf("Error sending greeting: %v", err)
		return
	}

	// Simple echo for demonstration (read client input and send it back)
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Printf("Error reading from client: %v", err)
			break
		}

		// Log client message
		clientMsg := string(buffer[:n])
		fmt.Printf("Received: %s", clientMsg)

		// For SMTP servers, you’d parse and respond based on the protocol,
		// but for now, we’ll just echo back the message for demo purposes.
		_, err = conn.Write([]byte("Echo: " + clientMsg))
		if err != nil {
			log.Printf("Error sending response: %v", err)
			break
		}
	}
}
