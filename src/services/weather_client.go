package services

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

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
