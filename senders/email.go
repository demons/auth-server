package senders

import (
	"fmt"
	"net/smtp"
)

// EmailConfig с настройками доступа к smtp-серверу
type EmailConfig struct {
	Host     string
	Port     int
	Login    string
	Password string
}

// EmailSender рассыльщик emil'ов
type EmailSender struct {
	config *EmailConfig
	from   string
}

// NewEmailSender создает EmailSender
func (c *EmailConfig) NewEmailSender(from string) *EmailSender {
	return &EmailSender{
		config: c,
		from:   from,
	}
}

// Send отправляет указанное письмо по электронной почте
func (s EmailSender) Send(to string, msg []byte) error {
	// Set up authentication information
	auth := smtp.PlainAuth("", s.config.Login, s.config.Password, s.config.Host)

	// Send email
	serverName := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	err := smtp.SendMail(serverName, auth, s.from, []string{to}, msg)

	return err
}
