package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"telemetry-service/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB represents the database connection
type DB struct {
	Pool *pgxpool.Pool
}

// New creates a new DB instance
func New(connString string) (*DB, error) {
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	// Test the connection
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return &DB{Pool: pool}, nil
}

// Close closes the database connection
func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

// UpdateTelemetry updates the telemetry value and status for a device
func (db *DB) UpdateTelemetry(ctx context.Context, deviceID int, value float64, status string) error {
	// First check if device exists
	var exists bool
	err := db.Pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM sensors WHERE id = $1)", deviceID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking device existence: %w", err)
	}
	if !exists {
		return errors.New("device not found")
	}

	// Update the sensor value and status
	query := `
		UPDATE sensors
		SET value = $1, status = $2, last_updated = $3
		WHERE id = $4
	`

	result, err := db.Pool.Exec(ctx, query, value, status, time.Now(), deviceID)
	if err != nil {
		return fmt.Errorf("error updating telemetry: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("device not found")
	}

	return nil
}

// GetTelemetry retrieves telemetry data for a specific device
func (db *DB) GetTelemetry(ctx context.Context, deviceID int) (models.TelemetryResponse, error) {
	query := `
		SELECT id, value, status, last_updated
		FROM sensors
		WHERE id = $1
	`

	var telemetry models.TelemetryResponse
	err := db.Pool.QueryRow(ctx, query, deviceID).Scan(
		&telemetry.DeviceID,
		&telemetry.Value,
		&telemetry.Status,
		&telemetry.LastUpdated,
	)
	if err != nil {
		return models.TelemetryResponse{}, fmt.Errorf("error getting telemetry: %w", err)
	}

	return telemetry, nil
}

// GetAllTelemetry retrieves telemetry data for all devices
func (db *DB) GetAllTelemetry(ctx context.Context) ([]models.TelemetryResponse, error) {
	query := `
		SELECT id, value, status, last_updated
		FROM sensors
		ORDER BY id
	`

	rows, err := db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying telemetry: %w", err)
	}
	defer rows.Close()

	var telemetries []models.TelemetryResponse
	for rows.Next() {
		var t models.TelemetryResponse
		err := rows.Scan(
			&t.DeviceID,
			&t.Value,
			&t.Status,
			&t.LastUpdated,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning telemetry row: %w", err)
		}
		telemetries = append(telemetries, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating telemetry rows: %w", err)
	}

	return telemetries, nil
}

