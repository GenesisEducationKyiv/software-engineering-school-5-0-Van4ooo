package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/services"
)

// nolint: goconst
func TestFetchCurrentWeather(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	apiKey := "1111"
	baseURL := "https://api.weatherapi.com/v1"
	expectedURL := fmt.Sprintf("%s/current.json?key=%s&q=Lviv", baseURL, apiKey)

	httpmock.RegisterResponder("GET", expectedURL,
		httpmock.NewStringResponder(
			http.StatusOK,
			`{"current": {"temp_c": 20.5, "humidity": 60, "condition": {"text": "Sunny"}}}`,
		),
	)

	cfg := &config.WeatherAPI{
		Key:     apiKey,
		BaseURL: baseURL,
	}

	weather, err := services.NewOpenWeatherService(cfg).GetWeather("Lviv")

	assert.NoError(t, err)
	assert.Equal(t, 20.5, weather.Temperature)
	assert.Equal(t, float64(60), weather.Humidity)
	assert.Equal(t, "Sunny", weather.Description)
}

func TestFetchCurrentWeatherError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	apiKey := "1111"
	baseURL := "https://api.weatherapi.com/v1"
	expectedURL := fmt.Sprintf("%s/current.json?key=%s&q=Lviv", baseURL, apiKey)

	httpmock.RegisterResponder("GET", expectedURL,
		httpmock.NewStringResponder(http.StatusNotFound, `{"error": "city not found"}`),
	)

	cfg := &config.WeatherAPI{
		Key:     apiKey,
		BaseURL: baseURL,
	}

	weather, err := services.NewOpenWeatherService(cfg).GetWeather("Lviv")

	assert.Error(t, err)
	assert.Nil(t, weather)
	assert.Contains(t, err.Error(), "city not found")
}
