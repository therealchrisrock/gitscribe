package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Database DatabaseConfig
	Firebase FirebaseConfig
	Server   ServerConfig
	User     UserConfig
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// FirebaseConfig holds Firebase configuration
type FirebaseConfig struct {
	ProjectID           string
	CredentialsPath     string
	UseEmulator         bool
	EmulatorHost        string
	ServiceAccountEmail string
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string
	Env  string
}

// UserConfig holds user module configuration
type UserConfig struct {
	RepositoryType string // "gorm" or "firebase"
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "your-super-secret-and-long-postgres-password"),
			Name:     getEnv("DB_NAME", "teammate_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Firebase: FirebaseConfig{
			ProjectID:           getEnv("FIREBASE_PROJECT_ID", ""),
			CredentialsPath:     getEnv("FIREBASE_CREDENTIALS_PATH", ""),
			UseEmulator:         getEnvBool("FIREBASE_USE_EMULATOR", false),
			EmulatorHost:        getEnv("FIREBASE_EMULATOR_HOST", "localhost:9099"),
			ServiceAccountEmail: getEnv("FIREBASE_SERVICE_ACCOUNT_EMAIL", ""),
		},
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Env:  getEnv("APP_ENV", "development"),
		},
		User: UserConfig{
			RepositoryType: getEnv("USER_REPOSITORY_TYPE", "gorm"),
		},
	}, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultValue
}

// getEnvBool gets an environment variable as boolean or returns a default value
func getEnvBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}
