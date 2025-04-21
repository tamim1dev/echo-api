// config/database.go
package config

import (
	"fmt"
	"os"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB initializes the database connection using Turso
func InitDB() (*gorm.DB, error) {
	// Get Turso database URL from environment variable
	dbURL := os.Getenv("TURSO_DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("TURSO_DATABASE_URL environment variable not set")
	}

	// Get Turso auth token from environment variable
	authToken := os.Getenv("TURSO_AUTH_TOKEN")
	if authToken == "" {
		return nil, fmt.Errorf("TURSO_AUTH_TOKEN environment variable not set")
	}

	// Configure SQLite driver with Turso options
	sqliteDialector := sqlite.Open(fmt.Sprintf("%s?_auth_token=%s", dbURL, authToken))

	// Configure GORM
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// Open connection
	db, err := gorm.Open(sqliteDialector, config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}
