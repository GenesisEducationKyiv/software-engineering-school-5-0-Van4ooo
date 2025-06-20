package email

import (
	"net/smtp"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/config"
)

type Template interface {
	GetTo() string
	GetMsg() []byte
}

func (s *Sender) Send(template Template) error {
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

type Sender struct {
	cfg  config.SMTPSettings
	auth smtp.Auth
}

func NewSender(cfg config.SMTPSettings) *Sender {
	return &Sender{
		cfg: cfg,
	}
}

func (s *Sender) ensureAuth() error {
	if s.auth == nil {
		s.auth = smtp.PlainAuth("", s.cfg.GetName(), s.cfg.GetPass(), s.cfg.GetHost())
	}
	return nil
}
