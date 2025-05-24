// @title CRUD Without DB API
// @version 1.0
// @description This is a sample server for a CRUD application without a database.
// @host localhost:3000
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
	router.Use(handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
	))

	// Add Swagger UI route
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

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
