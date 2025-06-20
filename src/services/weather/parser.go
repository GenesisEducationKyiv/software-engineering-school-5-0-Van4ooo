package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/models"
)

func Parser(data []byte, status int) (*models.Weather, error) {
	if status != http.StatusOK {
		return nil, errors.New("city not found")
	}
	var resp models.WeatherApiResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return &models.Weather{
		Temperature: resp.Current.TempC,
		Humidity:    float64(resp.Current.Humidity),
		Description: resp.Current.Condition.Text,
	}, nil
}

func FetchRaw(rawURL string) ([]byte, int, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid URL: %w", err)
	}

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Printf("warning: failed to close response body: %v", cerr)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	return body, resp.StatusCode, err
}
