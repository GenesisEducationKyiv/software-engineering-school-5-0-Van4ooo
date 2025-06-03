package services

import (
	"fmt"
	"io"
	"net/http"
)

func buildUrl(city, key string) string {
	return fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s", key, city)
}

func FetchRaw(url string) ([]byte, int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, 0, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	return body, resp.StatusCode, nil
}
