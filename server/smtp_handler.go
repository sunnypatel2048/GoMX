package server

import (
	"fmt"
	"net"
	"strings"
)

// Define the different states of an SMTP session
const (
	InitState = iota
	MailState
	RcptState
	DataState
	QuitState
)

// SMTPSession holds the state and details of an SMTP session.
type SMTPSession struct {
	State       int
	Sender      string
	Recipients  []string
	MessageBody strings.Builder
}

// NewSMTPSession initializes a new SMTP session with the initial state.
func NewSMTPSession() *SMTPSession {
	return &SMTPSession{
		State:      InitState,
		Recipients: make([]string, 0),
	}
}

// HandleSMTPCommands handles client commands during an SMTP session.
func HandleSMTPCommands(conn net.Conn) {
	defer conn.Close()

	session := NewSMTPSession()
	_, _ = conn.Write([]byte("220 Welcome to GoMX Mail Server\r\n"))

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Printf("Error reading from client: %v", err)
			break
		}

		// Parse and handle commands
		clientInput := string(buffer[:n])
		response := handleCommand(clientInput, session)
		if _, err = conn.Write([]byte(response)); err != nil {
			fmt.Printf("Error sending response: %v", err)
			break
		}

		if session.State == QuitState {
			break
		}
	}
}

// handleCommand parses and executes SMTP commands based on session state.
func handleCommand(command string, session *SMTPSession) string {
	command = strings.TrimSpace(command)
	parts := strings.SplitN(command, " ", 2)
	cmd := strings.ToUpper(parts[0])
	args := ""
	if len(parts) > 1 {
		args = parts[1]
	}

	switch cmd {
	case "HELO":
		if session.State == InitState {
			session.State = MailState
			return "250 Hello " + args + "\r\n"
		}
		return "503 Bad sequence of commands\r\n"

	case "MAIL":
		if strings.HasPrefix(args, "FROM:") && session.State == MailState {
			session.Sender = strings.TrimPrefix(args, "FROM:")
			session.State = RcptState
			return "250 Sender OK\r\n"
		}
		return "503 Bad sequence of commands\r\n"

	case "RCPT":
		if strings.HasPrefix(args, "TO:") && session.State == RcptState {
			recipient := strings.TrimPrefix(args, "TO:")
			session.Recipients = append(session.Recipients, recipient)
			return "250 Recipient OK\r\n"
		}
		return "503 Bad sequence of commands\r\n"

	case "DATA":
		if session.State == RcptState && len(session.Recipients) > 0 {
			session.State = DataState
			return "354 Start mail input; end with <CRLF>.<CRLF>\r\n"
		}
		return "503 Bad sequence of commands\r\n"

	case "QUIT":
		session.State = QuitState
		return "221 Bye\r\n"

	case "RSET":
		session.State = InitState
		session.Sender = ""
		session.Recipients = make([]string, 0)
		return "250 Reset OK\r\n"

	case "NOOP":
		return "250 OK\r\n"

	default:
		return "500 Command not recognized\r\n"
	}
}
