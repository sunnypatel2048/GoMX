package main

import (
	"GoMX/server"
	"fmt"
)

func main() {
	fmt.Println("Starting GoMX - a mail exchange server...")

	// Load configuration from config/config.json
	config := server.InitConfig()

	fmt.Printf("Loaded Configuration: Domain = %s, Port = %d\n", config.Domain, config.SMTPPort)

	// Start TCP server on the configured SMTP port
	go server.StartTCPServer(config.SMTPPort)

	// Block main thread to keep the server running
	select {}
}
