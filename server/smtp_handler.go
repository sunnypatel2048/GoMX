package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

// SessionState keeps track of the client's current state during an SMTP session.
type SessionState int

const (
	InitState     SessionState = iota // Initial state when the session starts.
	HeloState                         // After receiving a valid HELO command.
	MailFromState                     // After receiving a valid MAIL FROM command.
	RcptState                         // After receiving a valid RCPT TO command.
	DataState                         // After receiving a valid DATA command.
)

// SMTPSession stores session-specific information.
type SMTPSession struct {
	State       SessionState    // Current state of the SMTP session.
	Sender      string          // Sender's email address.
	Recipients  []string        // List of recipient email addresses.
	MessageBody strings.Builder // Stores the raw message body after the DATA command.
}

// handleSMTPConnection handles incoming SMTP client connections.
func handleSMTPConnection(conn net.Conn) {
	defer conn.Close()

	// Create a new SMTPSession for tracking the session state.
	session := &SMTPSession{State: InitState}

	// Send initial greeting to the client.
	sendResponse(conn, 220, "Welcome to GoMX Mail Server")

	reader := bufio.NewReader(conn)
	for {
		// Read the incoming command from the client.
		command, err := reader.ReadString('\n')
		if err != nil {
			sendResponse(conn, 421, "Service not available, closing transmission channel")
			return
		}

		// Trim command to remove extra whitespaces and newline characters.
		command = strings.TrimSpace(command)
		if len(command) == 0 {
			continue
		}

		// Split the command into command verb and arguments.
		parts := strings.SplitN(command, " ", 2)
		verb := strings.ToUpper(parts[0])
		var args string
		if len(parts) > 1 {
			args = parts[1]
		}

		// Handle the SMTP command based on the current state.
		response := handleCommand(session, verb, args)
		if response != "" {
			sendResponse(conn, parseSMTPCode(response), response)
		}

		// Close the connection if the client issues the QUIT command.
		if verb == "QUIT" {
			return
		}
	}
}

// handleCommand processes each SMTP command and updates the session state accordingly.
func handleCommand(session *SMTPSession, verb, args string) string {
	switch verb {
	case "HELO":
		if len(args) > 0 {
			session.State = HeloState
			return fmt.Sprintf("250 Hello %s", args)
		}
		return "501 Syntax error in parameters or arguments"

	case "MAIL":
		if strings.HasPrefix(strings.ToUpper(args), "FROM:") && session.State == HeloState {
			session.Sender = strings.TrimSpace(args[5:]) // Extract sender email.
			session.State = MailFromState
			return "250 Sender OK"
		}
		return getErrorForState(session.State, HeloState)

	case "RCPT":
		if strings.HasPrefix(strings.ToUpper(args), "TO:") && session.State == MailFromState {
			recipient := strings.TrimSpace(args[3:]) // Extract recipient email.
			session.Recipients = append(session.Recipients, recipient)
			session.State = RcptState
			return "250 Recipient OK"
		}
		return getErrorForState(session.State, MailFromState)

	case "DATA":
		if session.State == RcptState && len(session.Recipients) > 0 {
			session.State = DataState
			return "354 Start mail input; end with <CRLF>.<CRLF>"
		}
		return "503 Bad sequence of commands"

	case ".":
		if session.State == DataState {
			session.State = InitState
			return "250 Message accepted for delivery"
		}
		return "503 Bad sequence of commands"

	case "QUIT":
		return "221 Bye"

	default:
		return "500 Syntax error, command unrecognized"
	}
}

// parseSMTPCode extracts the SMTP response code from a string message.
func parseSMTPCode(message string) int {
	var code int
	_, err := fmt.Sscanf(message, "%d", &code)
	if err != nil {
		return 500 // Return 500 for unrecognized responses.
	}
	return code
}

// getErrorForState returns the appropriate error message for unexpected command sequences.
func getErrorForState(current, required SessionState) string {
	if current < required {
		return "503 Bad sequence of commands"
	}
	return "501 Syntax error in parameters or arguments"
}
