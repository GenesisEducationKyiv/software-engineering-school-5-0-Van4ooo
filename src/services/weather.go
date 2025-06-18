package services

import (
	"fmt"
	"log"
)

type WeatherService interface {
	GetWeather(city string) (*Weather, error)
}

type OpenWeatherService struct {
	APIKey  string
	BaseURL string
}

func NewOpenWeatherService(apiKey, baseURL string) *OpenWeatherService {
	return &OpenWeatherService{
		APIKey:  apiKey,
		BaseURL: baseURL,
	}
}

func (s *OpenWeatherService) GetWeather(city string) (*Weather, error) {
	log.Printf("Fetching weather for %q; APIKey=%s, BaseURL=%s",
		city, s.APIKey, s.BaseURL)

	url := fmt.Sprintf("%s/current.json?key=%s&q=%s", s.BaseURL, s.APIKey, city)
	raw, status, err := FetchRaw(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather: %w", err)
	}

	weather, err := ParserWeather(raw, status)
	if err != nil {
		return nil, fmt.Errorf("failed to parse weather data: %w", err)
	}
	return weather, nil
}
