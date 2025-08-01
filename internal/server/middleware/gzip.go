package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"
)

// GzipMiddleware compresses responses with gzip if supported by the browser
func GzipMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if client accepts gzip encoding
			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				next.ServeHTTP(w, r)
				return
			}

			// Create a wrapper to capture content type
			wrapper := &gzipWrapper{
				ResponseWriter: w,
				request:        r,
			}

			next.ServeHTTP(wrapper, r)
		})
	}
}

// gzipWrapper wraps http.ResponseWriter to conditionally enable gzip compression
type gzipWrapper struct {
	http.ResponseWriter
	request       *http.Request
	gzipWriter    *gzip.Writer
	headerWritten bool
}

func (w *gzipWrapper) WriteHeader(code int) {
	if w.headerWritten {
		return
	}
	w.headerWritten = true

	// Check if we should compress based on Content-Type
	contentType := w.Header().Get("Content-Type")
	if contentType == "" {
		// If no content type set, try to guess from path
		if ShouldCompress(w.request.URL.Path) {
			w.enableGzip()
		}
	} else {
		// Check if content type is compressible and not already compressed
		if ShouldCompressContentType(contentType) && !strings.Contains(contentType, "gzip") {
			w.enableGzip()
		}
	}

	w.ResponseWriter.WriteHeader(code)
}

func (w *gzipWrapper) Write(b []byte) (int, error) {
	if !w.headerWritten {
		w.WriteHeader(http.StatusOK)
	}

	if w.gzipWriter != nil {
		return w.gzipWriter.Write(b)
	}
	return w.ResponseWriter.Write(b)
}

func (w *gzipWrapper) enableGzip() {
	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Vary", "Accept-Encoding")
	w.gzipWriter = gzip.NewWriter(w.ResponseWriter)
}

// Close implements io.Closer for proper cleanup
func (w *gzipWrapper) Close() error {
	if w.gzipWriter != nil {
		return w.gzipWriter.Close()
	}
	return nil
}

// Flush implements http.Flusher
func (w *gzipWrapper) Flush() {
	if w.gzipWriter != nil {
		w.gzipWriter.Flush()
	}
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}
