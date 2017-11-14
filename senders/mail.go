package senders

import (
	"encoding/base64"
	"fmt"
)

// Mail информация и содержание письма
type Mail struct {
	Headers  map[string]string
	From     string
	FromDesc string
	Subject  string
	To       []string
	Body     string
}

// SetPlainText записывает в тело письма обычный текст
func (m *Mail) SetPlainText(text string) {
	headers := map[string]string{
		"MIME-version": "1.0",
		"Content-Type": "text/plane; charset=utf-8",
	}
	m.Headers = headers
	m.Body = text
}

// SetHTMLext записывает в тело письма html текст
func (m *Mail) SetHTMLext(html string) {
	headers := map[string]string{
		"MIME-version": "1.0",
		"Content-Type": "text/html; charset=utf-8",
	}
	m.Headers = headers
	m.Body = html
}

// Build выполняет сборку письма
func (m *Mail) Build() []byte {
	message := ""
	message += fmt.Sprintf("Subject: =?utf-8?b?%s?=\r\n", base64.StdEncoding.EncodeToString([]byte(m.Subject)))

	if m.FromDesc != "" && m.From != "" {
		message += fmt.Sprintf("From: =?utf-8?b?%s?= <%s>\r\n", base64.StdEncoding.EncodeToString([]byte(m.FromDesc)), m.From)
	} else if m.From != "" && m.FromDesc == "" {
		message += fmt.Sprintf("From: %s\r\n", m.From)
	}

	for key, value := range m.Headers {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	message += "\r\n" + m.Body

	return []byte(message)
}
