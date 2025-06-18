package config

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type ConfigSource interface {
	Load(provider EnvProvider) error
	Validate() error
}

type AppConfig struct {
	DB         Postgres
	SMTP       SMTP
	WeatherAPI WeatherAPI
}

type SMTP struct {
	Host string `validate:"required"`
	Addr string `validate:"required"`
	Name string `validate:"required"`
	Pass string `validate:"required"`
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

func (s *SMTP) Load(provider EnvProvider) error {
	s.Host = provider.Get("SMTP_HOST")
	s.Addr = provider.Get("SMTP_ADDR")
	s.Name = provider.Get("SMTP_NAME")
	s.Pass = provider.Get("SMTP_PASS")
	return nil
}

type WeatherAPI struct {
	Key     string `validate:"required"`
	BaseURL string `validate:"required"`
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

func (w *WeatherAPI) Load(provider EnvProvider) error {
	w.Key = provider.Get("WEATHER_API_KEY")
	w.BaseURL = provider.Get("WEATHER_API_BASE_URL")
	return nil
}

func (w *WeatherAPI) Validate() error {
	if err := validate.Struct(w); err != nil {
		return mapValidationErrorsToEnvVars(err, "WeatherAPI",
			map[string]string{
				"Key":     "WEATHER_API_KEY",
				"BaseURL": "WEATHER_API_BASE_URL",
			})
	}
	return nil
}

func mapValidationErrorsToEnvVars(err error, structName string,
	fieldToEnvNameMap map[string]string) error {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return err
	}

	var errorMessages []string
	for _, fieldErr := range validationErrors {
		envVarName, ok := fieldToEnvNameMap[fieldErr.Field()]
		if !ok {
			errorMessages = append(errorMessages, fmt.Sprintf("Field '%s' in %s: %s",
				fieldErr.Field(), structName, fieldErr.Tag()))
			continue
		}

		switch fieldErr.Tag() {
		case "required":
			errorMessages = append(errorMessages, fmt.Sprintf(
				"\n[!] Environment variable '%s' is missing or empty. "+
					"(Required for field '%s')",
				envVarName, fieldErr.Field()))
		case "url":
			errorMessages = append(errorMessages, fmt.Sprintf(
				"\n[!] Environment variable '%s' contains an invalid URL. "+
					"(Field '%s')",
				envVarName, fieldErr.Field()))
		default:
			errorMessages = append(errorMessages, fmt.Sprintf(
				"\n[!] Environment variable '%s' failed validation for tag '%s'. "+
					"(Field '%s')",
				envVarName, fieldErr.Tag(), fieldErr.Field()))
		}
	}

	if len(errorMessages) == 0 {
		return err
	}
	return errors.New(strings.Join(errorMessages, "; "))
}

func forEachConfigSourceField(targetStruct interface{},
	action func(loader ConfigSource, fieldName string) error) error {
	val := reflect.ValueOf(targetStruct).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := typ.Field(i).Name

		if field.CanInterface() {
			if loader, ok := field.Addr().Interface().(ConfigSource); ok {
				if err := action(loader, fieldName); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func NewAppConfig() *AppConfig {
	return &AppConfig{}
}

func (ac *AppConfig) Load(provider EnvProvider) error {
	ac.DB = Postgres{}
	ac.SMTP = SMTP{}
	ac.WeatherAPI = WeatherAPI{}

	return forEachConfigSourceField(ac, func(loader ConfigSource, fieldName string) error {
		if err := loader.Load(provider); err != nil {
			return fmt.Errorf("failed to load %s config: %w", fieldName, err)
		}
		return nil
	})
}

func (ac *AppConfig) Validate() error {
	return forEachConfigSourceField(ac, func(loader ConfigSource, fieldName string) error {
		if err := loader.Validate(); err != nil {
			return fmt.Errorf("%s config error: %w", fieldName, err)
		}
		return nil
	})
}

func Config() (*AppConfig, error) {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Warning: .env file not found or failed to load")
	}
	return SetupConfig(OSProvider{})
}

func SetupConfig(provider EnvProvider) (*AppConfig, error) {
	appConfig := NewAppConfig()

	if err := appConfig.Load(provider); err != nil {
		return nil, fmt.Errorf("failed to load application configuration: %w", err)
	}

	if err := appConfig.Validate(); err != nil {
		return nil, fmt.Errorf("application configuration validation failed: %w", err)
	}

	return appConfig, nil
}
