package services

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func buildUrl(city, key string) string {
	return fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s", key, city)
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
		if err := resp.Body.Close(); err != nil {
			log.Printf("warning: failed to close response body: %v", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	return body, resp.StatusCode, nil
}
