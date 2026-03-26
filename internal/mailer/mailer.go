package mailer

import (
	"emailsender/internal/config"
	"fmt"
	"mime"
	"net/smtp"
)

type Mailer struct {
	config config.Config
}

func New(configuration config.Config) *Mailer {
	return &Mailer{config: configuration}
}

func (mailer *Mailer) Send(to, subject, htmlBody string) error {
	auth := smtp.PlainAuth("", mailer.config.SMTPUser, mailer.config.SMTPPass, mailer.config.SMTPHost)

	encodedSubject := mime.QEncoding.Encode("utf-8", subject)
	encodedFromName := mime.QEncoding.Encode("utf-8", mailer.config.SMTPFromName)

	headers := fmt.Sprintf(
		"From: %s <%s>\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=utf-8\r\n\r\n",
		encodedFromName, mailer.config.SMTPFrom, to, encodedSubject,
	)

	message := []byte(headers + htmlBody)
	address := fmt.Sprintf("%s:%s", mailer.config.SMTPHost, mailer.config.SMTPPort)

	return smtp.SendMail(address, auth, mailer.config.SMTPFrom, []string{to}, message)
}
