package verification

import (
	"crypto/tls"
	"github.com/Inoi-K/Find-Me/pkg/config"
	"gopkg.in/gomail.v2"
)

type EmailData struct {
	Email   string
	URL     string
	Subject string
}

func send(data *EmailData) error {
	// prepare message
	m := gomail.NewMessage()

	m.SetHeader("From", config.C.EmailFrom)
	m.SetHeader("To", data.Email)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/plain", data.URL)

	d := gomail.NewDialer(config.C.SMTPHost, config.C.SMTPPort, config.C.SMTPUser, config.C.SMTPPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return d.DialAndSend(m)
}
