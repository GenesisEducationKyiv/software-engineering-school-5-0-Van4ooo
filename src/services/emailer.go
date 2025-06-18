package services

import (
	"fmt"
	"net/smtp"
)

func (s *SMTPEmailSender) Send(to, subject, body string) error {
	if err := s.ensureAuth(); err != nil {
		return err
	}
	msg := fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body)
	return smtp.SendMail(
		s.Addr,
		s.auth,
		s.Name,
		[]string{to},
		[]byte(msg),
	)
}

func (s *SMTPEmailSender) SendConfirmation(to, baseURL, token string) error {
	if err := s.ensureAuth(); err != nil {
		return err
	}
	link := fmt.Sprintf("%s/api/confirm/%s", baseURL, token)
	subject := "Confirm your subscription"
	message := fmt.Sprintf(
		"Subject: %s\r\n\r\nClick to confirm your subscription: %s", subject, link,
	)
	return smtp.SendMail(
		s.Addr,
		s.auth,
		s.Name,
		[]string{to},
		[]byte(message),
	)
}

type SMTPEmailSender struct {
	Host string
	Addr string
	Name string
	Pass string
	auth smtp.Auth
}

func NewSMTPEmailSender(host, addr, name, pass string) *SMTPEmailSender {
	return &SMTPEmailSender{
		Host: host,
		Addr: addr,
		Name: name,
		Pass: pass,
	}
}

func (s *SMTPEmailSender) ensureAuth() error {
	if s.auth == nil {
		s.auth = smtp.PlainAuth("", s.Name, s.Pass, s.Host)
	}
	return nil
}
