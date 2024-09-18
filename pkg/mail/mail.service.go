package mail

import (
	"fmt"
	"net/smtp"
)

type MailService struct {
	From     string
	Password string
	Host     string
	Port     uint
}

func (m *MailService) SendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", m.From, m.Password, m.Host)

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", m.From, to, subject, body)

	return smtp.SendMail(m.Host+":"+string(m.Port), auth, m.From, []string{to}, []byte(msg))
}
