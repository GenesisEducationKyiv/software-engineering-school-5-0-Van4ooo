package services

import (
	"fmt"
	"net/smtp"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/config"
)

type EmailSender interface {
	Send(to, subject, body string) error
	SendConfirmation(to, baseURL, token string) error
}

func (s *SMTPEmailSender) Send(to, subject, body string) error {
	if err := s.ensureAuth(); err != nil {
		return err
	}
	msg := fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body)
	return smtp.SendMail(
		s.cfg.GetAddr(),
		s.auth,
		s.cfg.GetName(),
		[]string{to},
		[]byte(msg),
	)
}

func (s *SMTPEmailSender) SendConfirmation(to, baseURL, token string) error {
	link := fmt.Sprintf("%s/api/confirm/%s", baseURL, token)
	subject := "Confirm your subscription"
	message := fmt.Sprintf(
		"Subject: %s\r\n\r\nClick to confirm your subscription: %s", subject, link,
	)
	return s.Send(to, subject, message)
}

type SMTPEmailSender struct {
	cfg  config.SMTPSettings
	auth smtp.Auth
}

func NewSMTPEmailSender(cfg config.SMTPSettings) *SMTPEmailSender {
	return &SMTPEmailSender{
		cfg: cfg,
	}
}

func (s *SMTPEmailSender) ensureAuth() error {
	if s.auth == nil {
		s.auth = smtp.PlainAuth("", s.cfg.GetName(), s.cfg.GetPass(), s.cfg.GetHost())
	}
	return nil
}
