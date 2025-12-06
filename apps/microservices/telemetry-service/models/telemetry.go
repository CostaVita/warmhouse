package models

import "time"

// Telemetry represents telemetry data for a device
type Telemetry struct {
	DeviceID    int       `json:"device_id"`
	Value       float64   `json:"value"`
	Status      string    `json:"status"`
	LastUpdated time.Time `json:"last_updated"`
}

// TelemetryUpdate represents the data to update telemetry
type TelemetryUpdate struct {
	Value  float64 `json:"value" binding:"required"`
	Status string  `json:"status" binding:"required"`
}

// TelemetryResponse represents the telemetry response
type TelemetryResponse struct {
	DeviceID    int       `json:"device_id"`
	Value       float64   `json:"value"`
	Status      string    `json:"status"`
	LastUpdated time.Time `json:"last_updated"`
}

