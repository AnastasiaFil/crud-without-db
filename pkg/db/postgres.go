package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
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
	// Try to load .env file (ignore error if file doesn't exist)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	port, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))

	config := &Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     port,
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", ""),
		DBName:   getEnv("DB_NAME", "crud_db"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	// Log configuration (without password for security)
	log.Printf("Database config: host=%s, port=%d, user=%s, dbname=%s, sslmode=%s",
		config.Host, config.Port, config.User, config.DBName, config.SSLMode)

	return config
}

func (c *Config) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

func NewPostgresConnection(config *Config) (*sql.DB, error) {
	log.Printf("Attempting to connect to PostgreSQL at %s:%d", config.Host, config.Port)

	db, err := sql.Open("postgres", config.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection with retry logic
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		if err := db.Ping(); err != nil {
			if i == maxRetries-1 {
				db.Close()
				return nil, fmt.Errorf("failed to ping database after %d attempts: %w", maxRetries, err)
			}
			log.Printf("Database connection attempt %d failed: %v. Retrying in 2 seconds...", i+1, err)
			time.Sleep(2 * time.Second)
			continue
		}
		break
	}

	log.Println("Successfully connected to PostgreSQL database")
	return db, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
