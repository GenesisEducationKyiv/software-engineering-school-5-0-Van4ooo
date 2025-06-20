package config

import (
	"fmt"
	"reflect"
)

type AppSettings interface {
	ConfigSource
	GetDB() DBSettings
	GetSMTP() SMTPSettings
	GetWeatherAPI() WeatherSettings
}

type AppConfig struct {
	DB         DBSettings
	SMTP       SMTPSettings
	WeatherAPI WeatherSettings
}

func NewAppConfig() *AppConfig {
	return &AppConfig{
		DB:         &Postgres{},
		SMTP:       &SMTP{},
		WeatherAPI: &WeatherAPI{},
	}
}

func (ac *AppConfig) Load(provider EnvProvider) error {
	val := reflect.ValueOf(ac).Elem()
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		src := val.Field(i).Interface().(ConfigSource)
		if err := src.Load(provider); err != nil {
			return fmt.Errorf("failed to load %s: %w", typ.Field(i).Name, err)
		}
	}
	return nil
}

func (ac *AppConfig) Validate() error {
	val := reflect.ValueOf(ac).Elem()
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		src := val.Field(i).Interface().(ConfigSource)
		if err := src.Validate(); err != nil {
			return fmt.Errorf("%s validation error: %w", typ.Field(i).Name, err)
		}
	}
	return nil
}

func (ac *AppConfig) GetDB() DBSettings              { return ac.DB }
func (ac *AppConfig) GetSMTP() SMTPSettings          { return ac.SMTP }
func (ac *AppConfig) GetWeatherAPI() WeatherSettings { return ac.WeatherAPI }

var _ AppSettings = (*AppConfig)(nil)
var _ ConfigSource = (*AppConfig)(nil)
