package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"telemetry-service/db"
	"telemetry-service/models"

	"github.com/gin-gonic/gin"
)

// TelemetryHandler handles telemetry-related requests
type TelemetryHandler struct {
	DB *db.DB
}

// NewTelemetryHandler creates a new TelemetryHandler
func NewTelemetryHandler(database *db.DB) *TelemetryHandler {
	return &TelemetryHandler{
		DB: database,
	}
}

// RegisterRoutes registers the telemetry routes
func (h *TelemetryHandler) RegisterRoutes(router *gin.RouterGroup) {
	telemetry := router.Group("/telemetry")
	{
		telemetry.GET("", h.GetAllTelemetry)
		telemetry.GET("/:deviceId", h.GetTelemetry)
		telemetry.PATCH("/:deviceId/value", h.UpdateTelemetry)
	}
}

// GetAllTelemetry handles GET /api/v1/telemetry
func (h *TelemetryHandler) GetAllTelemetry(c *gin.Context) {
	telemetries, err := h.DB.GetAllTelemetry(context.Background())
	if err != nil {
		log.Printf("Error getting all telemetry: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, telemetries)
}

// GetTelemetry handles GET /api/v1/telemetry/:deviceId
func (h *TelemetryHandler) GetTelemetry(c *gin.Context) {
	deviceID, err := strconv.Atoi(c.Param("deviceId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
		return
	}

	telemetry, err := h.DB.GetTelemetry(context.Background(), deviceID)
	if err != nil {
		log.Printf("Error getting telemetry: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Telemetry not found"})
		return
	}

	c.JSON(http.StatusOK, telemetry)
}

// UpdateTelemetry handles PATCH /api/v1/telemetry/:deviceId/value
func (h *TelemetryHandler) UpdateTelemetry(c *gin.Context) {
	deviceID, err := strconv.Atoi(c.Param("deviceId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
		return
	}

	var telemetryUpdate models.TelemetryUpdate
	if err := c.ShouldBindJSON(&telemetryUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.DB.UpdateTelemetry(context.Background(), deviceID, telemetryUpdate.Value, telemetryUpdate.Status)
	if err != nil {
		log.Printf("Error updating telemetry: %v", err)
		if err.Error() == "device not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Telemetry updated successfully"})
}

