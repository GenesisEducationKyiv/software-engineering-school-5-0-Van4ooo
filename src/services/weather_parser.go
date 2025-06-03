package services

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

type apiResponse struct {
	Current struct {
		TempC     float64 `json:"temp_c"`
		Humidity  int     `json:"humidity"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
}

type Weather struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Description string  `json:"description"`
}

func ParserWeather(data []byte, status int) (*Weather, error) {
	if status != http.StatusOK {
		return nil, errors.New("city not found")
	}

	var resp apiResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return &Weather{
		Temperature: resp.Current.TempC,
		Humidity:    float64(resp.Current.Humidity),
		Description: resp.Current.Condition.Text,
	}, nil
}

func FetchCurrentWeather(city string) (*Weather, error) {
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		return nil, errors.New("missing WEATHER_API_KEY")
	}

	raw, status, err := FetchRaw(buildUrl(city, apiKey))
	if err != nil {
		return nil, err
	}

	return ParserWeather(raw, status)
}
