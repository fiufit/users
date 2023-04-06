package utils

import (
	"context"
	"crypto/tls"
	"fmt"

	"firebase.google.com/go/v4/auth"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

const (
	verificationEmailBody = `<p>Hello,</p>
		<p>Follow <a href='%v'>this link</a> to verify your FiuFit email address.</p>
		<p>If you didn’t ask to verify this address, you can ignore this email.</p>
		<p>Thanks,</p>
		<p>Your the FiuFit team</p>`
)

type Mailer interface {
	SendEmail(userEmail, body string) error
	SendAccountVerificationEmail(ctx context.Context, email string) error
}

type MailerImpl struct {
	fromMail string
	password string
	smtpHost string
	smtpPort int
	auth     *auth.Client
	logger   *zap.Logger
}

func NewMailerImpl(fromMail, password, smtpHost string, smtpPort int, auth *auth.Client, logger *zap.Logger) *MailerImpl {
	return &MailerImpl{
		fromMail: fromMail,
		password: password,
		smtpHost: smtpHost,
		smtpPort: smtpPort,
		auth:     auth,
		logger:   logger,
	}
}

func (mailer MailerImpl) SendEmail(userEmail, body string) error {
	d := gomail.NewDialer(mailer.smtpHost, mailer.smtpPort, mailer.fromMail, mailer.password)
	msg := gomail.NewMessage()
	msg.SetBody("text/html", body)
	msg.SetHeader("To", userEmail)
	msg.SetHeader("From", mailer.fromMail)
	msg.SetHeader("Subject", "FiuFit Verification")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return d.DialAndSend(msg)
}

func (mailer MailerImpl) SendAccountVerificationEmail(ctx context.Context, email string) error {
	verificationLink, err := mailer.auth.EmailVerificationLink(ctx, email)
	if err != nil {
		mailer.logger.Error("Unable to generate verification link for email", zap.String("email", email), zap.Error(err))
		return err
	}

	/*TODO Find out how to user firebase's email verification instead of our own mail account + SMTP server. Apparently
	the Go firebase SDK doesn't have auth.SendEmailVerification()
	*/
	err = mailer.SendEmail(email, fmt.Sprintf(verificationEmailBody, verificationLink))
	if err != nil {
		return err
	}
	return nil
}
