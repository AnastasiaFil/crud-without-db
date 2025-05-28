package rest

import (
	"crud-without-db/pkg/logger"
	"net/http"
	"time"
)

// loggingMiddleware wraps an http.Handler to log the request details using zerolog
// It logs request method, URI, and response time in a structured format
func loggingMiddleware(next http.Handler) http.Handler {
	middlewareLogger := logger.GetLogger("middleware")

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Create a response writer wrapper to capture status code
			ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			// Process the request
			next.ServeHTTP(ww, r)

			// Log the request details
			duration := time.Since(start)
			middlewareLogger.Info().
				Str("method", r.Method).
				Str("uri", r.RequestURI).
				Str("remote_addr", r.RemoteAddr).
				Str("user_agent", r.UserAgent()).
				Int("status_code", ww.statusCode).
				Dur("duration", duration).
				Msg("HTTP request processed")
		},
	)
}

// responseWriter is a wrapper around http.ResponseWriter to capture the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code before writing it
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
