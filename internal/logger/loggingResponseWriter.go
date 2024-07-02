// Package sets Wraps a Response to use with logger to provide additional data on response
package logger

import (
	"net/http"
)

type (
	// responseData holds the response details for logging.
	responseData struct {
		size       int
		statusCode int
	}

	// logResponseWriter wraps http.ResponseWriter to capture response details.
	logResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

// Write captures the size of the response.
func (r *logResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

// WriteHeader captures the status code of the response.
func (r *logResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.statusCode = statusCode
}
