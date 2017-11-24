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
	config   *EmailConfig
	From     string
	FromDesc string
}

// NewEmailSender создает EmailSender
func (c *EmailConfig) NewEmailSender(from string, fromDesc string) *EmailSender {
	return &EmailSender{
		config:   c,
		From:     from,
		FromDesc: fromDesc,
	}
}

// Send отправляет указанное письмо по электронной почте
func (s EmailSender) Send(to string, msg []byte) error {
	// Set up authentication information
	auth := smtp.PlainAuth("", s.config.Login, s.config.Password, s.config.Host)

	// Send email
	serverName := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	err := smtp.SendMail(serverName, auth, s.From, []string{to}, msg)

	return err
}
