package database

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations executes database migrations
func RunMigrations(migrationsPath string) error {
	log.Printf("Running migrations from path: %s", migrationsPath)

	// Get the raw SQL database object from GORM
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database: %w", err)
	}

	return runMigrateInstance(sqlDB, migrationsPath)
}

// runMigrateInstance creates and runs a migrate instance
func runMigrateInstance(db *sql.DB, migrationsPath string) error {
	// Ensure path is absolute
	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Create the postgres driver instance
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	// Initialize the migrator
	sourceURL := fmt.Sprintf("file://%s", absPath)
	m, err := migrate.NewWithDatabaseInstance(
		sourceURL,
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	// Run the migrations
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	if err == migrate.ErrNoChange {
		log.Println("No migrations to run")
	} else {
		log.Println("Migrations completed successfully")
	}

	return nil
}

// CreateMigrationsTable ensures the migrations table exists
func CreateMigrationsTable() error {
	// Create the schema_migrations table if it doesn't exist
	query := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version bigint NOT NULL,
		dirty boolean NOT NULL,
		PRIMARY KEY (version)
	);`

	return DB.Exec(query).Error
}

// GetMigrationVersion returns the current migration version
func GetMigrationVersion() (int, bool, error) {
	// Check if table exists
	var exists bool
	err := DB.Raw(`SELECT EXISTS (
		SELECT FROM information_schema.tables 
		WHERE table_name = 'schema_migrations'
	)`).Scan(&exists).Error

	if err != nil {
		return 0, false, err
	}

	if !exists {
		return 0, false, nil
	}

	// Get migration version
	var version int
	var dirty bool
	err = DB.Raw(`SELECT version, dirty FROM schema_migrations LIMIT 1`).Row().Scan(&version, &dirty)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, false, nil
		}
		return 0, false, err
	}

	return version, dirty, nil
}
