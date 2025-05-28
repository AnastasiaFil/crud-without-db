package db

import (
	"crud-without-db/pkg/logger"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewConfigFromEnv() *Config {
	port, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))

	return &Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     port,
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", ""),
		DBName:   getEnv("DB_NAME", "postgres"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

func (c *Config) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

func NewPostgresConnection(config *Config) (*sql.DB, error) {
	dbLogger := logger.GetLogger("database")

	dbLogger.Info().
		Str("host", config.Host).
		Int("port", config.Port).
		Str("database", config.DBName).
		Str("user", config.User).
		Str("ssl_mode", config.SSLMode).
		Msg("Attempting to connect to PostgreSQL database")

	db, err := sql.Open("postgres", config.ConnectionString())
	if err != nil {
		dbLogger.Error().Err(err).Msg("Failed to open database connection")
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	dbLogger.Debug().
		Int("max_open_conns", 25).
		Int("max_idle_conns", 25).
		Dur("conn_max_lifetime", 5*time.Minute).
		Msg("Database connection pool configured")

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		dbLogger.Error().Err(err).Msg("Failed to ping database")
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	dbLogger.Info().Msg("Successfully connected to PostgreSQL database")
	return db, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
