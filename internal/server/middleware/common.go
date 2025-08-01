package middleware

import (
	"net/http"
	"strings"
)

// Common file extension sets
var (
	// StaticExtensions are common static file extensions
	StaticExtensions = []string{".css", ".js", ".png", ".jpg", ".jpeg", ".gif", ".svg", ".ico", ".woff", ".woff2", ".ttf", ".eot"}

	// CompressibleExtensions are file extensions that should be compressed
	CompressibleExtensions = []string{".html", ".css", ".js", ".json", ".xml", ".svg", ".txt"}

	// UncompressibleExtensions are file extensions that should not be compressed
	UncompressibleExtensions = []string{".png", ".jpg", ".jpeg", ".gif", ".webp", ".ico", ".mp4", ".webm", ".zip", ".gz", ".woff", ".woff2"}
)

// Middleware is a function that wraps an http.Handler
type Middleware func(http.Handler) http.Handler

// Chain applies multiple middleware functions to a handler
func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

// IsStaticFile checks if the path appears to be a static file
func IsStaticFile(path string) bool {
	for _, ext := range StaticExtensions {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}
	// Check for static paths
	return strings.HasPrefix(path, "/_/") || strings.HasPrefix(path, "/static/")
}

// ShouldCompress determines if the content should be compressed based on path
func ShouldCompress(path string) bool {
	// Don't compress already compressed files
	for _, ext := range UncompressibleExtensions {
		if strings.HasSuffix(path, ext) {
			return false
		}
	}

	// Compress text-based content
	for _, ext := range CompressibleExtensions {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}

	// Compress dynamic pages (no file extensions)
	return !strings.Contains(path, ".")
}

// ShouldCompressContentType determines if content should be compressed based on Content-Type
func ShouldCompressContentType(contentType string) bool {
	// Remove charset and other parameters
	ct := strings.Split(contentType, ";")[0]
	ct = strings.TrimSpace(ct)

	compressibleTypes := []string{
		"text/html",
		"text/css",
		"text/javascript",
		"text/plain",
		"application/javascript",
		"application/json",
		"application/xml",
		"image/svg+xml",
	}

	for _, compressible := range compressibleTypes {
		if ct == compressible {
			return true
		}
	}

	return false
}
