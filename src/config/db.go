package config

type DBSettings interface {
	ConfigSource
	GetURL() string
}

type Postgres struct {
	URL string `validate:"required"`
}

func (p *Postgres) Load(provider EnvProvider) error {
	p.URL = provider.Get("DATABASE_URL")
	return nil
}

func (p *Postgres) Validate() error {
	err := validate.Struct(p)
	if err != nil {
		return mapValidationErrorsToEnvVars(err, "Postgres",
			map[string]string{"URL": "DATABASE_URL"})
	}
	return nil
}

func (p *Postgres) GetURL() string { return p.URL }

var _ DBSettings = (*Postgres)(nil)
