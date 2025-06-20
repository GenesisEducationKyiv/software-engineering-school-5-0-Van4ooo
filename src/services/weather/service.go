package weather

import (
	"fmt"
	"log"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/models"
)

type Service struct {
	cfg config.WeatherSettings
}

func NewService(cfg config.WeatherSettings) *Service {
	return &Service{
		cfg: cfg,
	}
}

func (s *Service) GetByCity(city string) (*models.Weather, error) {
	log.Printf("Fetching weather for %q; APIKey=%s, BaseURL=%s",
		city, s.cfg.GetKey(), s.cfg.GetBaseURL())

	raw, status, err := FetchRaw(s.cfg.GenUrl(city))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather: %w", err)
	}

	weather, err := Parser(raw, status)
	if err != nil {
		return nil, fmt.Errorf("failed to parse weather data: %w", err)
	}
	return weather, nil
}
