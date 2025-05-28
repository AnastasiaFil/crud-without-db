package logger

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Config holds logger configuration
type Config struct {
	Level      string
	Format     string // "json" or "console"
	TimeFormat string
}

// NewConfigFromEnv creates logger config from environment variables
func NewConfigFromEnv() *Config {
	return &Config{
		Level:      getEnv("LOG_LEVEL", "info"),
		Format:     getEnv("LOG_FORMAT", "json"),
		TimeFormat: getEnv("LOG_TIME_FORMAT", time.RFC3339),
	}
}

// InitLogger initializes the global logger with the given configuration
func InitLogger(config *Config) {
	// Set log level
	level, err := zerolog.ParseLevel(config.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// Configure time format
	zerolog.TimeFieldFormat = config.TimeFormat

	// Configure output format
	if strings.ToLower(config.Format) == "console" {
		// Pretty console output for development
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "15:04:05",
		})
	} else {
		// JSON output for production
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}

	// Add caller information for debugging
	if level <= zerolog.DebugLevel {
		log.Logger = log.Logger.With().Caller().Logger()
	}
}

// GetLogger returns a logger with optional fields
func GetLogger(component string) zerolog.Logger {
	if component != "" {
		return log.With().Str("component", component).Logger()
	}
	return log.Logger
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
