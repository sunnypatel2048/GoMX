package main

import (
	"GoMX/auth"
	"GoMX/server"
	"GoMX/smtp"
	"GoMX/storage"
	"fmt"
)

func main() {
	fmt.Println("Start GoMX - a mail exchange server...")

	smtp.StartSMTPServer()
	auth.InitializeAuth()
	storage.InitStorage()
	server.StartServer()
}
