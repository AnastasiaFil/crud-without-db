// @title CRUD API with PostgreSQL
// @version 1.0
// @description This is a CRUD application with PostgreSQL database.
// @host
// @BasePath /
// @schemes http
package main

import (
	"context"
	_ "crud-without-db/docs"
	"crud-without-db/internal/repository/psql"
	"crud-without-db/internal/service"
	"crud-without-db/pkg/db"
	"crud-without-db/pkg/logger"
	"crud-without-db/pkg/rest"
	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Load .env file if it exists (optional for development)
	if err := godotenv.Load(); err != nil {
		// Don't fail if .env file doesn't exist - environment variables might be set another way
		log.Debug().Err(err).Msg("No .env file found, using system environment variables")
	}

	// Initialize logger
	logConfig := logger.NewConfigFromEnv()
	logger.InitLogger(logConfig)

	mainLogger := logger.GetLogger("main")
	mainLogger.Info().Msg("Starting CRUD API application")

	// Initialize database connection
	dbConfig := db.NewConfigFromEnv()
	database, err := db.NewPostgresConnection(dbConfig)
	if err != nil {
		mainLogger.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer database.Close()

	// Initialize repository with database connection
	usersRepo := psql.NewUsers(database)

	// Initialize database schema
	if err := usersRepo.InitSchema(); err != nil {
		mainLogger.Fatal().Err(err).Msg("Failed to initialize database schema")
	}

	// Initialize service and handler
	usersService := service.NewUsers(usersRepo)
	handler := rest.NewHandler(usersService)
	router := handler.InitRouter()

	// Enhanced CORS configuration for Swagger UI
	corsHandler := handlers.CORS(
		// Allow all origins for development - in production, specify your domain
		handlers.AllowedOrigins([]string{"*"}),
		// Allow all necessary HTTP methods
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"}),
		// Allow all necessary headers
		handlers.AllowedHeaders([]string{
			"Accept",
			"Accept-Language",
			"Content-Type",
			"Content-Language",
			"Origin",
			"Authorization",
			"X-Requested-With",
			"X-HTTP-Method-Override",
		}),
		// Allow credentials if needed
		handlers.AllowCredentials(),
		// Expose headers that might be needed
		handlers.ExposedHeaders([]string{"Content-Length", "Content-Type"}),
		// Cache preflight requests for 24 hours
		handlers.MaxAge(86400),
	)

	// Apply CORS middleware to the router
	router.Use(corsHandler)

	// Add Swagger UI route with custom configuration
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Add a health endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy", "database": "connected", "cors": "enabled"}`))
	}).Methods("GET", "OPTIONS")

	// Dynamic swagger.json endpoint that uses the current request host
	router.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Get the host from the request
		host := r.Host
		if host == "" {
			host = "localhost:3000" // fallback
		}

		// Create dynamic swagger JSON with current host
		swaggerJSON := `{
    "schemes": ["http"],
    "swagger": "2.0",
    "info": {
        "description": "This is a CRUD application with PostgreSQL database.",
        "title": "CRUD API with PostgreSQL",
        "contact": {},
        "version": "1.0"
    },
    "host": "` + host + `",
    "basePath": "/",
    "paths": {
        "/users": {
            "get": {
                "description": "Get a list of all users",
                "produces": ["application/json"],
                "tags": ["users"],
                "summary": "Get all users",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {"$ref": "#/definitions/domain.User"}
                        }
                    }
                }
            },
            "post": {
                "description": "Create a new user",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["users"],
                "summary": "Create a new user",
                "parameters": [{
                    "description": "Create user",
                    "name": "user",
                    "in": "body",
                    "required": true,
                    "schema": {"$ref": "#/definitions/domain.User"}
                }],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {"$ref": "#/definitions/domain.User"}
                    }
                }
            }
        },
        "/users/{id}": {
            "get": {
                "description": "Get a user by their ID",
                "produces": ["application/json"],
                "tags": ["users"],
                "summary": "Get a user by ID",
                "parameters": [{
                    "type": "integer",
                    "description": "User ID",
                    "name": "id",
                    "in": "path",
                    "required": true
                }],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {"$ref": "#/definitions/domain.User"}
                    }
                }
            },
            "put": {
                "description": "Update a user by their ID",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["users"],
                "summary": "Update a user",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Update user",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {"$ref": "#/definitions/domain.User"}
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {"$ref": "#/definitions/domain.User"}
                    }
                }
            },
            "delete": {
                "description": "Delete a user by their ID",
                "produces": ["application/json"],
                "tags": ["users"],
                "summary": "Delete a user",
                "parameters": [{
                    "type": "integer",
                    "description": "User ID",
                    "name": "id",
                    "in": "path",
                    "required": true
                }],
                "responses": {
                    "204": {"description": "No Content"}
                }
            }
        }
    },
    "definitions": {
        "domain.User": {
            "type": "object",
            "properties": {
                "age": {"type": "integer"},
                "id": {"type": "integer"},
                "name": {"type": "string"},
                "sex": {"type": "string"}
            }
        }
    }
}`

		w.Write([]byte(swaggerJSON))
	}).Methods("GET")

	// Initialize & run server
	srv := &http.Server{
		Addr:    ":3000",
		Handler: router,
	}

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start a goroutine to handle shutdown
	go func() {
		sig := <-sigChan
		mainLogger.Info().Str("signal", sig.String()).Msg("Received shutdown signal")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			mainLogger.Error().Err(err).Msg("Server shutdown failed")
		} else {
			mainLogger.Info().Msg("Server shut down gracefully")
		}
	}()

	// Start the server
	mainLogger.Info().
		Str("address", ":3000").
		Str("db_host", dbConfig.Host).
		Int("db_port", dbConfig.Port).
		Str("db_name", dbConfig.DBName).
		Msg("Server starting")

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		mainLogger.Fatal().Err(err).Msg("Failed to start server")
	}
}
