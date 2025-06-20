package email

import (
	"net/smtp"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/config"
)

type EmailTemplate interface {
	GetTo() string
	GetMsg() []byte
}

func (s *SMTPEmailSender) Send(template EmailTemplate) error {
	if err := s.ensureAuth(); err != nil {
		return err
	}

	return smtp.SendMail(
		s.cfg.GetAddr(),
		s.auth,
		s.cfg.GetName(),
		[]string{template.GetTo()},
		template.GetMsg(),
	)
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
