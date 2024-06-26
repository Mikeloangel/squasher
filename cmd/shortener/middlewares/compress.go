package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"slices"
	"strings"
)

type (
	compressWriter struct {
		http.ResponseWriter
		Writer io.Writer
	}
)

func (w compressWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func Compress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentTypes := []string{
			"application/javascript",
			"application/json",
			"text/css",
			"text/html",
			"text/plain",
			"text/xml",
		}
		acceptsEncoding := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
		compressableType := slices.Contains(contentTypes, r.Header.Get("Content-Type"))

		if !acceptsEncoding || !compressableType {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")

		cw := compressWriter{ResponseWriter: w, Writer: gz}
		next.ServeHTTP(cw, r)
	})
}
