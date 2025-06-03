package tests

import (
	"github.com/Van4ooo/genesis_case_task/src/services"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestFetchRaw(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	expectedURL := "https://api.weatherapi.com/v1/current.json?key=mock_api_key&q=London"
	httpmock.RegisterResponder("GET", expectedURL,
		httpmock.NewStringResponder(200, `{"location": {"name": "Lviv"}}`))

	body, statusCode, err := services.FetchRaw(expectedURL)

	assert.NoError(t, err)
	assert.Equal(t, 200, statusCode)
	assert.JSONEq(t, `{"location": {"name": "Lviv"}}`, string(body))
}

func TestFetchRawError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	expectedURL := "https://api.weatherapi.com/v1/current.json?key=mock_api_key&q=Lviv"
	httpmock.RegisterResponder("GET", expectedURL,
		httpmock.NewStringResponder(404, `{"error": "city not found"}`))

	body, statusCode, err := services.FetchRaw(expectedURL)

	assert.NoError(t, err)
	assert.Equal(t, 404, statusCode)
	assert.Equal(t, []byte(`{"error": "city not found"}`), body)
}

func TestParserWeather(t *testing.T) {
	response := []byte(`{"current": {"temp_c": 20.5, "humidity": 60, "condition": {"text": "Sunny"}}}`)
	status := 200

	weather, err := services.ParserWeather(response, status)

	assert.NoError(t, err)
	assert.Equal(t, 20.5, weather.Temperature)
	assert.Equal(t, float64(60), weather.Humidity)
	assert.Equal(t, "Sunny", weather.Description)
}

func TestParserWeatherError(t *testing.T) {
	response := []byte(`{"error": "city not found"}`)
	status := 404

	weather, err := services.ParserWeather(response, status)

	assert.Error(t, err)
	assert.Nil(t, weather)
}

func TestFetchCurrentWeather(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	apiKey := os.Getenv("WEATHER_API_KEY")
	expectedURL := "https://api.weatherapi.com/v1/current.json?key=" + apiKey + "&q=Lviv"

	httpmock.RegisterResponder("GET", expectedURL,
		httpmock.NewStringResponder(200, `{"current": {"temp_c": 20.5, "humidity": 60, "condition": {"text": "Sunny"}}}`))

	weather, err := services.FetchCurrentWeather("Lviv")

	assert.NoError(t, err)
	assert.Equal(t, 20.5, weather.Temperature)
	assert.Equal(t, float64(60), weather.Humidity)
	assert.Equal(t, "Sunny", weather.Description)
}

func TestFetchCurrentWeatherError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	apiKey := os.Getenv("WEATHER_API_KEY")
	expectedURL := "https://api.weatherapi.com/v1/current.json?key=" + apiKey + "&q=Lviv"

	httpmock.RegisterResponder("GET", expectedURL,
		httpmock.NewStringResponder(404, `{"error": "city not found"}`))

	weather, err := services.FetchCurrentWeather("Lviv")

	assert.Error(t, err)
	assert.Equal(t, "city not found", err.Error())
	assert.Nil(t, weather)
}
