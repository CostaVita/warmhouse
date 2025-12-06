package models

import "time"

// DeviceType represents the type of device
type DeviceType string

const (
	Temperature DeviceType = "temperature"
)

// Device represents a smart home device/sensor
type Device struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Type        DeviceType `json:"type"`
	Location    string     `json:"location"`
	Unit        string     `json:"unit"`
	Status      string     `json:"status"`
	LastUpdated time.Time  `json:"last_updated"`
	CreatedAt   time.Time  `json:"created_at"`
}

// DeviceCreate represents the data needed to create a new device
type DeviceCreate struct {
	Name     string     `json:"name" binding:"required"`
	Type     DeviceType `json:"type" binding:"required"`
	Location string     `json:"location" binding:"required"`
	Unit     string     `json:"unit"`
}

// DeviceUpdate represents the data that can be updated for a device
type DeviceUpdate struct {
	Name     string     `json:"name"`
	Type     DeviceType `json:"type"`
	Location string     `json:"location"`
	Unit     string     `json:"unit"`
	Status   string     `json:"status"`
}

