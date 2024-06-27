package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"slices"
	"strings"
)

type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (cw compressWriter) Write(b []byte) (int, error) {
	return cw.zw.Write(b)
}

func (cw compressWriter) Header() http.Header {
	return cw.w.Header()
}

func (cw compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		cw.w.Header().Set("Content-Encoding", "gzip")
	}
	cw.w.WriteHeader(statusCode)
}

func (cw compressWriter) Close() error {
	return cw.zw.Close()
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (cr compressReader) Read(p []byte) (n int, err error) {
	return cr.zr.Read(p)
}

func (cr compressReader) Close() error {
	if err := cr.r.Close(); err != nil {
		return err
	}

	return cr.zr.Close()
}

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

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
		useCompress := acceptsEncoding && compressableType

		if useCompress {
			cw := newCompressWriter(w)
			ow = cw
			defer cw.Close()
		}

		contentEncoding := strings.Contains(r.Header.Get("Content-Encoding"), "gzip")
		useDecompress := contentEncoding && compressableType

		if useDecompress {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close()
		}

		next.ServeHTTP(ow, r)
	})
}
