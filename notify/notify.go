package notify

import (
	"bytes"
	"html/template"
	"log"

	"audiolang.com/auth-server/senders"
)

// EmailNotificator занимается уведомлением пользователя по электронной почте
type EmailNotificator struct {
	sender *senders.EmailSender
}

// NewEmailNotificator создает новый email notificator
func NewEmailNotificator(sender *senders.EmailSender) *EmailNotificator {
	return &EmailNotificator{
		sender: sender,
	}
}

// SendActivationCode отправить код активации пользователю
func (n EmailNotificator) SendActivationCode(template *template.Template, to string, code string) {
	if template == nil {
		log.Println("Error sending email. Template is nil")
		return
	}

	data := struct {
		Link string
	}{
		Link: "http://audiolang.com/account/activate/?code=" + code,
	}

	// Компилируем шаблон
	buf := new(bytes.Buffer)
	if err := template.Execute(buf, data); err != nil {
		log.Printf("Error compiling template: %v", err)
		return
	}

	// Создаем письмо
	mail := senders.Mail{
		From:     "notify@audiolang.com",
		FromDesc: "Audiolang",
		Subject:  "Активация аккаунта",
	}
	mail.SetHTMLext(buf.String())

	// Компилируем и отправляем письмо
	if err := n.sender.Send(to, mail.Build()); err != nil {
		log.Printf("Error sending email: %v", err)
		return
	}
}

// SendResetPasswordMessage отправляет письмо для восстановления пароля
func (n EmailNotificator) SendResetPasswordMessage(template *template.Template, to string, code string) {
	if template == nil {
		log.Println("Error sending email. Template is nil")
		return
	}

	data := struct {
		Link string
	}{
		Link: "http://audiolang.com/account/password/change/?code=" + code,
	}

	// Компилируем шаблон
	buf := new(bytes.Buffer)
	if err := template.Execute(buf, data); err != nil {
		log.Printf("Error compiling template: %v", err)
		return
	}

	// Создаем письмо
	mail := senders.Mail{
		From:     "notify@audiolang.com",
		FromDesc: "Audiolang",
		Subject:  "Восстановление пароля",
	}
	mail.SetHTMLext(buf.String())

	// Компилируем и отправляем письмо
	if err := n.sender.Send(to, mail.Build()); err != nil {
		log.Printf("Error sending email: %v", err)
		return
	}
}
