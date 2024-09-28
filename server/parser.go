package server

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

// EmailMessage represents a parsed email message according to RFC 5322.
type EmailMessage struct {
	Headers map[string]string
	Body    string
}

// ParseEmailMessage parses a raw email message into its headers and body.
func ParseEmailMessage(rawData string) (*EmailMessage, error) {
	message := &EmailMessage{Headers: make(map[string]string)}

	// Use a scanner to read the raw message line by line.
	scanner := bufio.NewScanner(strings.NewReader(rawData))
	isHeader := true
	var previousHeader string

	for scanner.Scan() {
		line := scanner.Text()

		// If we encounter an empty line, switch to parsing the body.
		if line == "" {
			isHeader = false
			continue
		}

		if isHeader {
			// Headers may be continued on the next line if they start with a space or tab.
			if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") {
				message.Headers[previousHeader] += " " + strings.TrimSpace(line)
			} else {
				// Extract header name and value using regex.
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					headerName := strings.TrimSpace(parts[0])
					headerValue := strings.TrimSpace(parts[1])
					message.Headers[headerName] = headerValue
					previousHeader = headerName
				}
			}
		} else {
			// If we are in the body section, accumulate the body lines.
			message.Body += line + "\n"
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning message: %v", err)
	}

	return message, nil
}

// ValidateHeaders validates the basic structure of the required email headers.
func ValidateHeaders(headers map[string]string) error {
	requiredHeaders := []string{"From", "To", "Date"}
	for _, header := range requiredHeaders {
		if _, exists := headers[header]; !exists {
			return fmt.Errorf("missing required header: %s", header)
		}
	}

	// Validate "From" header.
	from := headers["From"]
	if !isValidEmail(from) {
		return fmt.Errorf("invalid From header format: %s", from)
	}

	// Validate "To" header.
	to := headers["To"]
	if !isValidEmail(to) {
		return fmt.Errorf("invalid To header format: %s", to)
	}

	return nil
}

// isValidEmail uses a simple regex to validate email addresses.
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
