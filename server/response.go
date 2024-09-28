package server

import (
	"fmt"
	"net"
)

// sendResponse sends a formatted SMTP response to the client.
func sendResponse(conn net.Conn, code int, message string) {
	response := fmt.Sprintf("%d %s\r\n", code, message)
	fmt.Fprintf(conn, response)
}
