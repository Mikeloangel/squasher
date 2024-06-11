package middlewares

import "net/http"

func TextPlain(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			contentType := r.Header.Get("Content-Type")
			if contentType != "text/plain" {
				http.Error(w, "Content type mismatch", http.StatusBadRequest)
				return
			}
			next.ServeHTTP(w, r)
		},
	)
}
