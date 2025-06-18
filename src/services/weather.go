package services

import (
	"fmt"
	"log"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/config"
)

type WeatherService interface {
	GetWeather(city string) (*Weather, error)
}

type OpenWeatherService struct {
	cfg config.WeatherSettings
}

func NewOpenWeatherService(cfg config.WeatherSettings) *OpenWeatherService {
	return &OpenWeatherService{
		cfg: cfg,
	}
}

func (s *OpenWeatherService) GetWeather(city string) (*Weather, error) {
	log.Printf("Fetching weather for %q; APIKey=%s, BaseURL=%s",
		city, s.cfg.GetKey(), s.cfg.GetBaseURL())

	raw, status, err := FetchRaw(s.cfg.GenUrl(city))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather: %w", err)
	}

	weather, err := ParserWeather(raw, status)
	if err != nil {
		return nil, fmt.Errorf("failed to parse weather data: %w", err)
	}
	return weather, nil
}
