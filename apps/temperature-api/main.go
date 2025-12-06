package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type TemperatureResponse struct {
	Value       float64   `json:"value"`
	Unit        string    `json:"unit"`
	Timestamp   time.Time `json:"timestamp"`
	Location    string    `json:"location"`
	Status      string    `json:"status"`
	SensorID    string    `json:"sensor_id"`
	SensorType  string    `json:"sensor_type"`
	Description string    `json:"description"`
}

func main() {
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Temperature endpoint with location query parameter
	router.GET("/temperature", getTemperatureByLocation)

	// Temperature endpoint with sensor ID
	router.GET("/temperature/:sensorID", getTemperatureBySensorID)

	// Get port from environment or use default
	port := getEnv("PORT", "8081")
	addr := ":" + port

	log.Printf("Temperature API server starting on %s\n", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}

// getTemperatureByLocation handles GET /temperature?location=
func getTemperatureByLocation(c *gin.Context) {
	location := c.Query("location")
	sensorID := c.Query("sensor_id")

	// If no location is provided, use a default based on sensor ID
	if location == "" {
		switch sensorID {
		case "1":
			location = "Living Room"
		case "2":
			location = "Bedroom"
		case "3":
			location = "Kitchen"
		default:
			location = "Unknown"
		}
	}

	// If no sensor ID is provided, generate one based on location
	if sensorID == "" {
		switch location {
		case "Living Room":
			sensorID = "1"
		case "Bedroom":
			sensorID = "2"
		case "Kitchen":
			sensorID = "3"
		default:
			sensorID = "0"
		}
	}

	// Generate random temperature between 18.0 and 25.0 degrees Celsius
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	temperature := 18.0 + r.Float64()*(25.0-18.0)
	// Round to 1 decimal place
	temperature = float64(int(temperature*10+0.5)) / 10

	response := TemperatureResponse{
		Value:       temperature,
		Unit:        "Celsius",
		Timestamp:   time.Now(),
		Location:    location,
		Status:      "active",
		SensorID:    sensorID,
		SensorType:  "temperature",
		Description: "Temperature sensor reading",
	}

	c.JSON(http.StatusOK, response)
}

// getTemperatureBySensorID handles GET /temperature/:sensorID
func getTemperatureBySensorID(c *gin.Context) {
	sensorID := c.Param("sensorID")
	location := c.Query("location")

	// If no location is provided, use a default based on sensor ID
	if location == "" {
		switch sensorID {
		case "1":
			location = "Living Room"
		case "2":
			location = "Bedroom"
		case "3":
			location = "Kitchen"
		default:
			location = "Unknown"
		}
	}

	// Generate random temperature between 18.0 and 25.0 degrees Celsius
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	temperature := 18.0 + r.Float64()*(25.0-18.0)
	// Round to 1 decimal place
	temperature = float64(int(temperature*10+0.5)) / 10

	response := TemperatureResponse{
		Value:       temperature,
		Unit:        "Celsius",
		Timestamp:   time.Now(),
		Location:    location,
		Status:      "active",
		SensorID:    sensorID,
		SensorType:  "temperature",
		Description: "Temperature sensor reading",
	}

	c.JSON(http.StatusOK, response)
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

