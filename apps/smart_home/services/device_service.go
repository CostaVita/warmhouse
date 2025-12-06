package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// DeviceService handles communication with device-registry service
type DeviceService struct {
	BaseURL    string
	HTTPClient *http.Client
}

// Device represents a device from device-registry service
type Device struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Location    string    `json:"location"`
	Unit        string    `json:"unit"`
	Status      string    `json:"status"`
	LastUpdated time.Time `json:"last_updated"`
	CreatedAt   time.Time `json:"created_at"`
}

// DeviceCreate represents the data needed to create a device
type DeviceCreate struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Location string `json:"location"`
	Unit     string `json:"unit"`
}

// DeviceUpdate represents the data that can be updated for a device
type DeviceUpdate struct {
	Name     string  `json:"name,omitempty"`
	Type     string  `json:"type,omitempty"`
	Location string  `json:"location,omitempty"`
	Unit     string  `json:"unit,omitempty"`
	Status   string  `json:"status,omitempty"`
	Value    *float64 `json:"value,omitempty"`
}

// NewDeviceService creates a new device service
func NewDeviceService(baseURL string) *DeviceService {
	return &DeviceService{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetDevices fetches all devices
func (s *DeviceService) GetDevices() ([]Device, error) {
	url := fmt.Sprintf("%s/api/v1/devices", s.BaseURL)

	resp, err := s.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching devices: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var devices []Device
	if err := json.NewDecoder(resp.Body).Decode(&devices); err != nil {
		return nil, fmt.Errorf("error decoding devices response: %w", err)
	}

	return devices, nil
}

// GetDeviceByID fetches a device by ID
func (s *DeviceService) GetDeviceByID(id int) (*Device, error) {
	url := fmt.Sprintf("%s/api/v1/devices/%d", s.BaseURL, id)

	resp, err := s.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching device: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var device Device
	if err := json.NewDecoder(resp.Body).Decode(&device); err != nil {
		return nil, fmt.Errorf("error decoding device response: %w", err)
	}

	return &device, nil
}

// CreateDevice creates a new device
func (s *DeviceService) CreateDevice(deviceCreate DeviceCreate) (*Device, error) {
	url := fmt.Sprintf("%s/api/v1/devices", s.BaseURL)

	jsonData, err := json.Marshal(deviceCreate)
	if err != nil {
		return nil, fmt.Errorf("error marshaling device create: %w", err)
	}

	resp, err := s.HTTPClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating device: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var device Device
	if err := json.NewDecoder(resp.Body).Decode(&device); err != nil {
		return nil, fmt.Errorf("error decoding device response: %w", err)
	}

	return &device, nil
}

// UpdateDevice updates an existing device
func (s *DeviceService) UpdateDevice(id int, deviceUpdate DeviceUpdate) (*Device, error) {
	url := fmt.Sprintf("%s/api/v1/devices/%d", s.BaseURL, id)

	jsonData, err := json.Marshal(deviceUpdate)
	if err != nil {
		return nil, fmt.Errorf("error marshaling device update: %w", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error updating device: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var device Device
	if err := json.NewDecoder(resp.Body).Decode(&device); err != nil {
		return nil, fmt.Errorf("error decoding device response: %w", err)
	}

	return &device, nil
}

// DeleteDevice deletes a device
func (s *DeviceService) DeleteDevice(id int) error {
	url := fmt.Sprintf("%s/api/v1/devices/%d", s.BaseURL, id)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error deleting device: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

