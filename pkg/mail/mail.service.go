package mail

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendEmail(to, subject, body string) error {
	from := os.Getenv("MAIL_FROM")
	password := os.Getenv("MAIL_PASSWORD")
	host := os.Getenv("MAIL_HOST")
	port := os.Getenv("MAIL_PORT")
	auth := smtp.PlainAuth("", from, password, host)

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, to, subject, body)

	return smtp.SendMail(host+":"+port, auth, from, []string{to}, []byte(msg))
}
