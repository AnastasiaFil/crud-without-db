package rest

import (
	"log"
	"net/http"
	"time"
)

// loggingMiddleware wraps an http.Handler to log the request details (timestamp, method, and URI)
// before delegating to the original handler. It returns a new handler and does not modify the input.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%s: [%s] - %s ", time.Now().Format(time.RFC3339), r.Method, r.RequestURI)
			next.ServeHTTP(w, r)
		},
	)
}
