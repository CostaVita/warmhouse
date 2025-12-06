package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"smarthome/models"
	"smarthome/services"

	"github.com/gin-gonic/gin"
)

// SensorHandler handles sensor-related requests
type SensorHandler struct {
	DeviceService      *services.DeviceService
	TelemetryService   *services.TelemetryService
	TemperatureService *services.TemperatureService
}

// NewSensorHandler creates a new SensorHandler
func NewSensorHandler(
	deviceService *services.DeviceService,
	telemetryService *services.TelemetryService,
	temperatureService *services.TemperatureService,
) *SensorHandler {
	return &SensorHandler{
		DeviceService:      deviceService,
		TelemetryService:   telemetryService,
		TemperatureService: temperatureService,
	}
}

// RegisterRoutes registers the sensor routes
func (h *SensorHandler) RegisterRoutes(router *gin.RouterGroup) {
	sensors := router.Group("/sensors")
	{
		sensors.GET("", h.GetSensors)
		sensors.GET("/:id", h.GetSensorByID)
		sensors.POST("", h.CreateSensor)
		sensors.PUT("/:id", h.UpdateSensor)
		sensors.DELETE("/:id", h.DeleteSensor)
		sensors.PATCH("/:id/value", h.UpdateSensorValue)
		sensors.GET("/temperature/:location", h.GetTemperatureByLocation)
	}
}

// GetSensors handles GET /api/v1/sensors
func (h *SensorHandler) GetSensors(c *gin.Context) {
	// Get devices from device-registry
	devices, err := h.DeviceService.GetDevices()
	if err != nil {
		log.Printf("Error getting devices from device-registry: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get telemetry for all devices and combine with device data
	var sensors []models.Sensor
	for _, device := range devices {
		sensor := models.Sensor{
			ID:          device.ID,
			Name:        device.Name,
			Type:        models.SensorType(device.Type),
			Location:    device.Location,
			Unit:        device.Unit,
			Status:      device.Status,
			LastUpdated: device.LastUpdated,
			CreatedAt:   device.CreatedAt,
			Value:       0, // Default value
		}

		// Get telemetry data for this device
		telemetry, err := h.TelemetryService.GetTelemetry(device.ID)
		if err == nil {
			sensor.Value = telemetry.Value
			sensor.Status = telemetry.Status
			sensor.LastUpdated = telemetry.LastUpdated
		}

		// Update temperature sensors with real-time data from the external API
		if sensor.Type == models.Temperature {
			tempData, err := h.TemperatureService.GetTemperatureByID(fmt.Sprintf("%d", sensor.ID))
			if err == nil {
				// Update sensor with real-time data
				sensor.Value = tempData.Value
				sensor.Status = tempData.Status
				sensor.LastUpdated = tempData.Timestamp
				log.Printf("Updated temperature data for sensor %d from external API", sensor.ID)
			} else {
				log.Printf("Failed to fetch temperature data for sensor %d: %v", sensor.ID, err)
			}
		}

		sensors = append(sensors, sensor)
	}

	c.JSON(http.StatusOK, sensors)
}

// GetSensorByID handles GET /api/v1/sensors/:id
func (h *SensorHandler) GetSensorByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sensor ID"})
		return
	}

	// Get device from device-registry
	device, err := h.DeviceService.GetDeviceByID(id)
	if err != nil {
		log.Printf("Error getting device from device-registry: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Sensor not found"})
		return
	}

	sensor := models.Sensor{
		ID:          device.ID,
		Name:        device.Name,
		Type:        models.SensorType(device.Type),
		Location:    device.Location,
		Unit:        device.Unit,
		Status:      device.Status,
		LastUpdated: device.LastUpdated,
		CreatedAt:   device.CreatedAt,
		Value:       0, // Default value
	}

	// Get telemetry data for this device
	telemetry, err := h.TelemetryService.GetTelemetry(id)
	if err == nil {
		sensor.Value = telemetry.Value
		sensor.Status = telemetry.Status
		sensor.LastUpdated = telemetry.LastUpdated
	}

	// If this is a temperature sensor, fetch real-time data from the temperature API
	if sensor.Type == models.Temperature {
		tempData, err := h.TemperatureService.GetTemperatureByID(fmt.Sprintf("%d", sensor.ID))
		if err == nil {
			// Update sensor with real-time data
			sensor.Value = tempData.Value
			sensor.Status = tempData.Status
			sensor.LastUpdated = tempData.Timestamp
			log.Printf("Updated temperature data for sensor %d from external API", sensor.ID)
		} else {
			log.Printf("Failed to fetch temperature data for sensor %d: %v", sensor.ID, err)
		}
	}

	c.JSON(http.StatusOK, sensor)
}

// GetTemperatureByLocation handles GET /api/v1/sensors/temperature/:location
func (h *SensorHandler) GetTemperatureByLocation(c *gin.Context) {
	location := c.Param("location")
	if location == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Location is required"})
		return
	}

	// Fetch temperature data from the external API
	tempData, err := h.TemperatureService.GetTemperature(location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to fetch temperature data: %v", err),
		})
		return
	}

	// Return the temperature data
	c.JSON(http.StatusOK, gin.H{
		"location":    tempData.Location,
		"value":       tempData.Value,
		"unit":        tempData.Unit,
		"status":      tempData.Status,
		"timestamp":   tempData.Timestamp,
		"description": tempData.Description,
	})
}

// CreateSensor handles POST /api/v1/sensors
func (h *SensorHandler) CreateSensor(c *gin.Context) {
	var sensorCreate models.SensorCreate
	if err := c.ShouldBindJSON(&sensorCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create device via device-registry
	deviceCreate := services.DeviceCreate{
		Name:     sensorCreate.Name,
		Type:     string(sensorCreate.Type),
		Location: sensorCreate.Location,
		Unit:     sensorCreate.Unit,
	}

	device, err := h.DeviceService.CreateDevice(deviceCreate)
	if err != nil {
		log.Printf("Error creating device via device-registry: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert device to sensor format
	sensor := models.Sensor{
		ID:          device.ID,
		Name:        device.Name,
		Type:        models.SensorType(device.Type),
		Location:    device.Location,
		Unit:        device.Unit,
		Status:      device.Status,
		LastUpdated: device.LastUpdated,
		CreatedAt:   device.CreatedAt,
		Value:       0,
	}

	c.JSON(http.StatusCreated, sensor)
}

// UpdateSensor handles PUT /api/v1/sensors/:id
func (h *SensorHandler) UpdateSensor(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sensor ID"})
		return
	}

	var sensorUpdate models.SensorUpdate
	if err := c.ShouldBindJSON(&sensorUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update device via device-registry
	deviceUpdate := services.DeviceUpdate{
		Name:     sensorUpdate.Name,
		Location: sensorUpdate.Location,
		Unit:     sensorUpdate.Unit,
		Status:   sensorUpdate.Status,
		Value:    sensorUpdate.Value,
	}

	if sensorUpdate.Type != "" {
		deviceUpdate.Type = string(sensorUpdate.Type)
	}

	device, err := h.DeviceService.UpdateDevice(id, deviceUpdate)
	if err != nil {
		log.Printf("Error updating device via device-registry: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert device to sensor format
	sensor := models.Sensor{
		ID:          device.ID,
		Name:        device.Name,
		Type:        models.SensorType(device.Type),
		Location:    device.Location,
		Unit:        device.Unit,
		Status:      device.Status,
		LastUpdated: device.LastUpdated,
		CreatedAt:   device.CreatedAt,
		Value:       0,
	}

	// Get telemetry if available
	telemetry, err := h.TelemetryService.GetTelemetry(id)
	if err == nil {
		sensor.Value = telemetry.Value
	}

	c.JSON(http.StatusOK, sensor)
}

// DeleteSensor handles DELETE /api/v1/sensors/:id
func (h *SensorHandler) DeleteSensor(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sensor ID"})
		return
	}

	// Delete device via device-registry
	err = h.DeviceService.DeleteDevice(id)
	if err != nil {
		log.Printf("Error deleting device via device-registry: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sensor deleted successfully"})
}

// UpdateSensorValue handles PATCH /api/v1/sensors/:id/value
func (h *SensorHandler) UpdateSensorValue(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sensor ID"})
		return
	}

	var request struct {
		Value  float64 `json:"value" binding:"required"`
		Status string  `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update telemetry via telemetry-service
	err = h.TelemetryService.UpdateTelemetry(id, request.Value, request.Status)
	if err != nil {
		log.Printf("Error updating telemetry via telemetry-service: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sensor value updated successfully"})
}
