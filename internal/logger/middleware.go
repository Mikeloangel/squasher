// Pakage handles middlewares for logger
package logger

import (
	"net/http"
	"time"
)

// WithLoggerMiddleware logs details about the HTTP request and response.
func WithLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uri := r.URL.Path
		method := r.Method
		start := time.Now()

		loggingWriter := wrapResponseWriter(w)
		next.ServeHTTP(loggingWriter, r)
		duration := time.Since(start)

		Log.Sugar().Infoln(
			"uri", uri,
			"method", method,
			"duration", duration,
			"status", loggingWriter.responseData.statusCode,
			"size", loggingWriter.responseData.size,
		)
	})
}

// wrapResponseWriter wraps Response writer with injected reponseData to be captured
func wrapResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	responseData := &responseData{
		size:       0,
		statusCode: 0,
	}

	return &loggingResponseWriter{
		ResponseWriter: w,
		responseData:   responseData,
	}
}
