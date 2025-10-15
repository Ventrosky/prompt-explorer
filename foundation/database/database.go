package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/ventrosky/prompt-explorer/foundation/config"
)

// Database represents a database connection
type Database struct {
	db        *sql.DB
	config    *config.Config
	connected bool
}

// New creates a new Database instance
func New(cfg *config.Config) *Database {
	return &Database{
		config:    cfg,
		connected: false,
	}
}

// Connect establishes a database connection
func (d *Database) Connect() error {
	if d.connected {
		return nil
	}

	// Build connection string
	connStr := d.config.DatabaseURL
	if connStr == "" {
		connStr = "postgres://localhost/prompt_explorer?sslmode=disable"
	}

	// Create database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		db.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}

	d.db = db
	d.connected = true
	log.Println("Database connected successfully")

	return nil
}

// Close closes the database connection
func (d *Database) Close() {
	if d.db != nil {
		d.db.Close()
		d.connected = false
		log.Println("Database connection closed")
	}
}

// IsConnected returns true if database is connected
func (d *Database) IsConnected() bool {
	return d.connected
}

// Ping checks if database is reachable
func (d *Database) Ping() error {
	if !d.connected || d.db == nil {
		return fmt.Errorf("database not connected")
	}

	return d.db.Ping()
}

// GetDB returns the database connection
func (d *Database) GetDB() *sql.DB {
	return d.db
}
