package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// TelemetryService handles communication with telemetry-service
type TelemetryService struct {
	BaseURL    string
	HTTPClient *http.Client
}

// TelemetryResponse represents telemetry data
type TelemetryResponse struct {
	DeviceID    int       `json:"device_id"`
	Value       float64   `json:"value"`
	Status      string    `json:"status"`
	LastUpdated time.Time `json:"last_updated"`
}

// TelemetryUpdate represents the data to update telemetry
type TelemetryUpdate struct {
	Value  float64 `json:"value"`
	Status string  `json:"status"`
}

// NewTelemetryService creates a new telemetry service
func NewTelemetryService(baseURL string) *TelemetryService {
	return &TelemetryService{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// UpdateTelemetry updates telemetry value for a device
func (s *TelemetryService) UpdateTelemetry(deviceID int, value float64, status string) error {
	url := fmt.Sprintf("%s/api/v1/telemetry/%d/value", s.BaseURL, deviceID)

	update := TelemetryUpdate{
		Value:  value,
		Status: status,
	}

	jsonData, err := json.Marshal(update)
	if err != nil {
		return fmt.Errorf("error marshaling telemetry update: %w", err)
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error updating telemetry: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetTelemetry fetches telemetry for a specific device
func (s *TelemetryService) GetTelemetry(deviceID int) (*TelemetryResponse, error) {
	url := fmt.Sprintf("%s/api/v1/telemetry/%d", s.BaseURL, deviceID)

	resp, err := s.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching telemetry: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var telemetry TelemetryResponse
	if err := json.NewDecoder(resp.Body).Decode(&telemetry); err != nil {
		return nil, fmt.Errorf("error decoding telemetry response: %w", err)
	}

	return &telemetry, nil
}

