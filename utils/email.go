package utils

import (
	"crypto/tls"

	"gopkg.in/gomail.v2"
)

type Mailer interface {
	SendAccountVerificationEmail(userEmail, verificationLink string) error
}

type MailerImpl struct {
	fromMail string
	password string
	smtpHost string
	smtpPort int
}

func NewMailerImpl(fromMail, password, smtpHost string, smtpPort int) *MailerImpl {
	return &MailerImpl{
		fromMail: fromMail,
		password: password,
		smtpHost: smtpHost,
		smtpPort: smtpPort,
	}
}

func (mailer MailerImpl) SendAccountVerificationEmail(userEmail, verificationLink string) error {
	d := gomail.NewDialer(mailer.smtpHost, mailer.smtpPort, mailer.fromMail, mailer.password)
	msg := gomail.NewMessage()
	msg.SetBody("text/html", "Gracias por registrarte con fiufit che, acá tenés el link de verificación: "+verificationLink)
	msg.SetHeader("To", userEmail)
	msg.SetHeader("From", mailer.fromMail)
	msg.SetHeader("Subject", "FiuFit Verification")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return d.DialAndSend(msg)
}
