package mailer

import (
	"crypto/tls"
	"emailsender/internal/config"
	"fmt"
	"mime"
	"net"
	"net/smtp"
	"net/url"
	"time"

	"golang.org/x/net/proxy"
)

const sendTimeout = 30 * time.Second

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

	connection, error := mailer.dial(address)
	if error != nil {
		return error
	}
	defer connection.Close()
	connection.SetDeadline(time.Now().Add(sendTimeout))

	client, error := smtp.NewClient(connection, mailer.config.SMTPHost)
	if error != nil {
		return error
	}
	defer client.Close()

	if ok, _ := client.Extension("STARTTLS"); ok {
		if error := client.StartTLS(&tls.Config{ServerName: mailer.config.SMTPHost}); error != nil {
			return error
		}
	}
	if error := client.Auth(auth); error != nil {
		return error
	}
	if error := client.Mail(mailer.config.SMTPFrom); error != nil {
		return error
	}
	if error := client.Rcpt(to); error != nil {
		return error
	}
	writer, error := client.Data()
	if error != nil {
		return error
	}
	if _, error := writer.Write(message); error != nil {
		return error
	}
	if error := writer.Close(); error != nil {
		return error
	}
	return client.Quit()
}

func (mailer *Mailer) dial(address string) (net.Conn, error) {
	directDialer := &net.Dialer{Timeout: sendTimeout}
	if mailer.config.SMTPProxy == "" {
		return directDialer.Dial("tcp", address)
	}

	proxyURL, error := url.Parse(mailer.config.SMTPProxy)
	if error != nil {
		return nil, fmt.Errorf("invalid SMTP_PROXY: %w", error)
	}
	proxyDialer, error := proxy.FromURL(proxyURL, directDialer)
	if error != nil {
		return nil, fmt.Errorf("invalid SMTP_PROXY: %w", error)
	}
	return proxyDialer.Dial("tcp", address)
}
