package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"weather-api/internal/models"
)

func GetWeatherData(ctx context.Context) (models.WeatherResponse, error) {
	var data models.WeatherResponse

	rawURL := os.Getenv("URL")
	if strings.TrimSpace(rawURL) == "" {
		return data, fmt.Errorf("URL env var not set")
	}
	key := os.Getenv("KEY")
	if strings.TrimSpace(key) == "" {
		return data, fmt.Errorf("KEY env var not set")
	}

	callURL, err := buildCallURL(rawURL, key)
	if err != nil {
		return data, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, callURL, nil)
	if err != nil {
		return data, err
	}

	resp, err := http.DefaultClient.Do(req)
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

func buildCallURL(rawURL, key string) (string, error) {
	if strings.Contains(rawURL, "{KEY}") {
		return strings.Replace(rawURL, "{KEY}", url.QueryEscape(key), 1), nil
	}
	if strings.Contains(rawURL, "%s") {
		return strings.Replace(rawURL, "%s", url.QueryEscape(key), 1), nil
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	q := parsed.Query()
	if q.Get("key") == "" {
		q.Set("key", key)
		parsed.RawQuery = q.Encode()
	}
	return parsed.String(), nil
}
