package config

import (
	"errors"
	"fmt"
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
