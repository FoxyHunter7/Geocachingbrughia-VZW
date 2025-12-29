package middleware

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"net/http"
)

// responseWriter wraps http.ResponseWriter to capture the response for ETag generation
type etagResponseWriter struct {
	http.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (w *etagResponseWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

func (w *etagResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

// CacheControl adds ETag-based caching headers.
// Uses "no-cache" which means the browser must always revalidate with the server,
// but can use the cached response if ETag matches (304 Not Modified).
// This ensures content updates are immediately visible while still benefiting from caching.
func CacheControl() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Capture the response body for ETag generation
			ew := &etagResponseWriter{
				ResponseWriter: w,
				body:           &bytes.Buffer{},
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(ew, r)

			// Generate ETag from response body
			hash := md5.Sum(ew.body.Bytes())
			etag := `"` + hex.EncodeToString(hash[:]) + `"`

			// Check if client sent If-None-Match header
			if match := r.Header.Get("If-None-Match"); match == etag {
				w.Header().Set("ETag", etag)
				w.Header().Set("Cache-Control", "no-cache")
				w.WriteHeader(http.StatusNotModified)
				return
			}

			// Set cache headers - no-cache means always revalidate, but can use cache if ETag matches
			w.Header().Set("ETag", etag)
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Vary", "Accept-Language, Accept-Encoding")

			// Write the actual response
			w.WriteHeader(ew.statusCode)
			w.Write(ew.body.Bytes())
		})
	}
}

// NoCache prevents caching for sensitive endpoints
func NoCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, private")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		next.ServeHTTP(w, r)
	})
}
