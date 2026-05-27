package core_mailer

import (
	"fmt"
	"net/smtp"
)

type Mailer interface {
	SendEmail(to, subject, body string) error
}

type SMTPMailer struct {
	config Config
}

func NewSMTPMailer(config Config) *SMTPMailer {
	return &SMTPMailer{config: config}
}

func (m *SMTPMailer) SendEmail(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", m.config.Host, m.config.Port)
	auth := smtp.PlainAuth("", m.config.Username, m.config.Password, m.config.Host)

	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=utf-8\r\n\r\n%s",
		m.config.From, to, subject, body,
	)

	return smtp.SendMail(addr, auth, m.config.From, []string{to}, []byte(msg))
}
