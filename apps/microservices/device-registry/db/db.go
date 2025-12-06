package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"device-registry/models"

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

// GetDevices retrieves all devices from the database
func (db *DB) GetDevices(ctx context.Context) ([]models.Device, error) {
	query := `
		SELECT id, name, type, location, unit, status, last_updated, created_at
		FROM sensors
		ORDER BY id
	`

	rows, err := db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying devices: %w", err)
	}
	defer rows.Close()

	var devices []models.Device
	for rows.Next() {
		var d models.Device
		err := rows.Scan(
			&d.ID,
			&d.Name,
			&d.Type,
			&d.Location,
			&d.Unit,
			&d.Status,
			&d.LastUpdated,
			&d.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning device row: %w", err)
		}
		devices = append(devices, d)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating device rows: %w", err)
	}

	return devices, nil
}

// GetDeviceByID retrieves a device by its ID
func (db *DB) GetDeviceByID(ctx context.Context, id int) (models.Device, error) {
	query := `
		SELECT id, name, type, location, unit, status, last_updated, created_at
		FROM sensors
		WHERE id = $1
	`

	var d models.Device
	err := db.Pool.QueryRow(ctx, query, id).Scan(
		&d.ID,
		&d.Name,
		&d.Type,
		&d.Location,
		&d.Unit,
		&d.Status,
		&d.LastUpdated,
		&d.CreatedAt,
	)
	if err != nil {
		return models.Device{}, fmt.Errorf("error getting device by ID: %w", err)
	}

	return d, nil
}

// CreateDevice creates a new device in the database
func (db *DB) CreateDevice(ctx context.Context, d models.DeviceCreate) (models.Device, error) {
	query := `
		INSERT INTO sensors (name, type, location, unit, status, last_updated, created_at)
		VALUES ($1, $2, $3, $4, 'inactive', $5, $5)
		RETURNING id, name, type, location, unit, status, last_updated, created_at
	`

	now := time.Now()
	var device models.Device
	err := db.Pool.QueryRow(ctx, query,
		d.Name,
		d.Type,
		d.Location,
		d.Unit,
		now,
	).Scan(
		&device.ID,
		&device.Name,
		&device.Type,
		&device.Location,
		&device.Unit,
		&device.Status,
		&device.LastUpdated,
		&device.CreatedAt,
	)
	if err != nil {
		return models.Device{}, fmt.Errorf("error creating device: %w", err)
	}

	return device, nil
}

// UpdateDevice updates an existing device
func (db *DB) UpdateDevice(ctx context.Context, id int, d models.DeviceUpdate) (models.Device, error) {
	// First check if the device exists
	_, err := db.GetDeviceByID(ctx, id)
	if err != nil {
		return models.Device{}, err
	}

	// Build the update query dynamically based on which fields are provided
	query := "UPDATE sensors SET last_updated = $1"
	args := []interface{}{time.Now()}
	argCount := 2

	if d.Name != "" {
		query += fmt.Sprintf(", name = $%d", argCount)
		args = append(args, d.Name)
		argCount++
	}

	if d.Type != "" {
		query += fmt.Sprintf(", type = $%d", argCount)
		args = append(args, d.Type)
		argCount++
	}

	if d.Location != "" {
		query += fmt.Sprintf(", location = $%d", argCount)
		args = append(args, d.Location)
		argCount++
	}

	if d.Unit != "" {
		query += fmt.Sprintf(", unit = $%d", argCount)
		args = append(args, d.Unit)
		argCount++
	}

	if d.Status != "" {
		query += fmt.Sprintf(", status = $%d", argCount)
		args = append(args, d.Status)
		argCount++
	}

	// Add the WHERE clause and RETURNING clause
	query += ` WHERE id = $` + fmt.Sprintf("%d", argCount) + `
		RETURNING id, name, type, location, unit, status, last_updated, created_at`
	args = append(args, id)

	var device models.Device
	err = db.Pool.QueryRow(ctx, query, args...).Scan(
		&device.ID,
		&device.Name,
		&device.Type,
		&device.Location,
		&device.Unit,
		&device.Status,
		&device.LastUpdated,
		&device.CreatedAt,
	)
	if err != nil {
		return models.Device{}, fmt.Errorf("error updating device: %w", err)
	}

	return device, nil
}

// DeleteDevice deletes a device by its ID
func (db *DB) DeleteDevice(ctx context.Context, id int) error {
	query := "DELETE FROM sensors WHERE id = $1"
	result, err := db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting device: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("device not found")
	}

	return nil
}

