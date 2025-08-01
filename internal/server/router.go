package server

import (
	"compress/gzip"
	"log/slog"
	"net/http"
	"strings"

	"github.com/lumiknit/can/internal/config"
	"github.com/lumiknit/can/internal/server/middleware"
	"github.com/lumiknit/can/pkg/css"
	"github.com/lumiknit/can/pkg/web"
	"github.com/lumiknit/can/public"
)

func SetupRoutes(logger *slog.Logger, cfg *config.Config) http.Handler {
	mux := http.NewServeMux()

	// Mount public files at root path
	mountPublicFiles(mux, logger, cfg.BasePath)

	// Mount all apps
	mountAllApps(mux, cfg.BasePath)

	// Health check endpoint
	mux.HandleFunc(cfg.BasePath+"healthz", healthzHandler)

	// Apply middleware
	handler := middleware.Chain(
		mux,
		// middleware.GzipMiddleware(),
		middleware.CommonHeadersMiddleware(cfg.AllowedOrigins()),
		middleware.RecoveryMiddleware(logger),
		middleware.LoggingMiddleware(logger),
	)

	return handler
}

// mountPublicFiles mounts public static files at root path
func mountPublicFiles(mux *http.ServeMux, logger *slog.Logger, basePath string) {
	staticFiles, err := web.EmbedFSToStaticFiles(public.PublicFS, public.PublicRoot)
	if err != nil {
		logger.Error("Failed to load static files", "error", err)
		panic(err)
	}

	for _, file := range staticFiles {
		data := file.Content
		contentType := file.ContentType

		// Process content based on type
		switch strings.Split(contentType, ";")[0] {
		case "text/css":
			if stylesheet, err := css.Parse(string(file.Content)); err == nil {
				data = []byte(stylesheet.Minify())
			} else {
				logger.Warn("Failed to parse CSS for minification", "path", file.Href, "error", err)
			}
		}

		// Pre-compress compressible static files
		var gzipData []byte
		if middleware.ShouldCompressContentType(contentType) {
			if compressed, err := precompressGzip(data); err == nil {
				gzipData = compressed
			} else {
				logger.Warn("Failed to pre-compress static file", "path", file.Href, "error", err)
			}
		}

		// Mount at configured base path (remove /_/ prefix)
		rootPath := strings.TrimPrefix(file.Href, "/_")
		if rootPath == "" {
			rootPath = "/"
		}

		// Add base path prefix
		if basePath != "/" {
			rootPath = strings.TrimSuffix(basePath, "/") + rootPath
		}

		// Capture variables for closure
		originalContent := data
		compressedContent := gzipData
		mux.HandleFunc(rootPath, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", contentType)

			if r.Method == http.MethodHead {
				w.WriteHeader(http.StatusOK)
				return
			}

			// Use pre-compressed version if client supports gzip
			if compressedContent != nil && strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				w.Header().Set("Content-Encoding", "gzip")
				w.Header().Set("Vary", "Accept-Encoding")
				w.Write(compressedContent)
			} else {
				w.Write(originalContent)
			}
		})
	}
}

// precompressGzip compresses data with gzip and returns the compressed bytes
func precompressGzip(data []byte) ([]byte, error) {
	var buf strings.Builder
	gz := gzip.NewWriter(&buf)

	if _, err := gz.Write(data); err != nil {
		gz.Close()
		return nil, err
	}

	if err := gz.Close(); err != nil {
		return nil, err
	}

	return []byte(buf.String()), nil
}

// mountApp mounts a web app to the mux
func mountApp(mux *http.ServeMux, app *web.App) {
	// Mount static files with app's base path prefix
	for _, file := range app.StaticFiles {
		path := app.BasePath + strings.TrimPrefix(file.Href, "/_")
		if path == app.BasePath {
			path = app.BasePath + "/"
		}

		// Capture variables for closure
		contentType := file.ContentType
		content := file.Content
		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", contentType)
			if r.Method == http.MethodHead {
				w.WriteHeader(http.StatusOK)
				return
			}
			w.Write(content)
		})
	}

	// Mount page handlers
	for _, page := range app.Pages {
		pagePath := app.BasePath + strings.TrimPrefix(page.Path, "/")
		if pagePath == app.BasePath {
			pagePath = app.BasePath
		}
		if !strings.HasSuffix(pagePath, "/") && pagePath != app.BasePath {
			pagePath += "/"
		}

		mux.HandleFunc(pagePath, func(w http.ResponseWriter, r *http.Request) {
			app.HandlePage(w, r)
		})
	}
}
