package helpers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	"weather-api/internal/models"
)

var ErrLocationNotSet = errors.New("location not set")

func GetWeatherData(ctx context.Context, location string) (models.WeatherResponse, error) {
	var data models.WeatherResponse

	rawURL := os.Getenv("URL")
	if strings.TrimSpace(rawURL) == "" {
		return data, fmt.Errorf("URL env var not set")
	}
	key := os.Getenv("KEY")
	if strings.TrimSpace(key) == "" {
		return data, fmt.Errorf("KEY env var not set")
	}

	callURL, err := buildCallURL(rawURL, key, location)
	if err != nil {
		return data, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, callURL, nil)
	if err != nil {
		return data, err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return data, fmt.Errorf("weather api error: %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return data, err
	}

	return data, nil
}

func buildCallURL(rawURL, key, location string) (string, error) {
	hasLocationPlaceholder := strings.Contains(rawURL, "{LOCATION}")
	rawURL = strings.Replace(rawURL, "{KEY}", url.QueryEscape(key), 1)
	rawURL = strings.Replace(rawURL, "%s", url.QueryEscape(key), 1)
	rawURL = strings.Replace(rawURL, "{LOCATION}", url.PathEscape(location), 1)
	if hasLocationPlaceholder && strings.TrimSpace(location) == "" {
		return "", ErrLocationNotSet
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	if strings.TrimSpace(location) != "" && !hasLocationPlaceholder {
		if idx := strings.Index(parsed.Path, "/timeline/"); idx != -1 {
			base := parsed.Path[:idx+len("/timeline/")]
			parsed.Path = base + url.PathEscape(location)
		}
	}

	q := parsed.Query()
	if q.Get("key") == "" {
		q.Set("key", key)
		parsed.RawQuery = q.Encode()
	}
	return parsed.String(), nil
}
