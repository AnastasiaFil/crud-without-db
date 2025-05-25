// @title CRUD Without DB API
// @version 1.0
// @description This is a sample server for a CRUD application without a database.
// @host 16.171.25.228:3000
// @BasePath /
// @schemes http
package main

import (
	"context"
	_ "crud-without-db/docs"
	"crud-without-db/internal/repository/psql"
	"crud-without-db/internal/service"
	"crud-without-db/pkg/rest"
	"github.com/gorilla/handlers"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// init deps
	usersRepo := psql.NewUsers()
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

	// Add Swagger UI route
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Add a test endpoint to verify CORS is working
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy", "cors": "enabled"}`))
	}).Methods("GET", "OPTIONS")

	// init & run server
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
		log.Printf("Received signal: %v. Shutting down.", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown failed: %v", err)
		} else {
			log.Println("Server shut down gracefully")
		}
	}()

	// Start the server
	log.Println("Server starting on :3000")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}
