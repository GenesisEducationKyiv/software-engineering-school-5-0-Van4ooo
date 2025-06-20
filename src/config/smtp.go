package config

type SMTPSettings interface {
	ConfigSource
	GetHost() string
	GetAddr() string
	GetName() string
	GetPass() string
}

type SMTP struct {
	Host string `validate:"required"`
	Addr string `validate:"required"`
	Name string `validate:"required"`
	Pass string `validate:"required"`
}

func (s *SMTP) Load(provider EnvProvider) error {
	s.Host = provider.Get("SMTP_HOST")
	s.Addr = provider.Get("SMTP_ADDR")
	s.Name = provider.Get("SMTP_NAME")
	s.Pass = provider.Get("SMTP_PASS")
	return nil
}

func (s *SMTP) Validate() error {
	err := validate.Struct(s)
	if err != nil {
		return mapValidationErrorsToEnvVars(err, "SMTP", map[string]string{
			"Host": "SMTP_HOST",
			"Addr": "SMTP_ADDR",
			"Name": "SMTP_NAME",
			"Pass": "SMTP_PASS",
		})
	}
	return nil
}

func (s *SMTP) GetHost() string { return s.Host }
func (s *SMTP) GetAddr() string { return s.Addr }
func (s *SMTP) GetName() string { return s.Name }
func (s *SMTP) GetPass() string { return s.Pass }

var _ SMTPSettings = (*SMTP)(nil)
