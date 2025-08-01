package server

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/lumiknit/can/pkg/css"
	"github.com/lumiknit/can/pkg/web"
	"github.com/lumiknit/can/public"
)

func SetupRoutes(logger *slog.Logger, releaseMode bool) http.Handler {
	mux := http.NewServeMux()

	// Mount public files using EmbedFSToStaticFiles
	staticFiles, err := web.EmbedFSToStaticFiles(public.PublicFS, public.PublicRoot)
	if err != nil {
		logger.Error("Failed to load static files", "error", err)
		panic(err)
	}

	for _, file := range staticFiles {
		data := file.Content
		switch strings.Split(file.ContentType, ";")[0] {
		case "text/css":
			if stylesheet, err := css.Parse(string(file.Content)); err == nil {
				data = []byte(stylesheet.Minify())
			} else {
				logger.Warn("Failed to parse CSS for minification", "path", file.Href, "error", err)
			}
		}

		// Capture variables for closure
		contentType := file.ContentType
		content := data
		mux.HandleFunc(file.Href, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", contentType)
			if r.Method == http.MethodHead {
				w.WriteHeader(http.StatusOK)
				return
			}
			w.Write(content)
		})
	}

	// Health check endpoint
	mux.HandleFunc("/healthz", healthzHandler)

	// Index page
	mux.HandleFunc("/", indexHandler(logger))

	// Apply middleware
	handler := Chain(
		mux,
		RecoveryMiddleware(logger),
		LoggingMiddleware(logger),
	)

	return handler
}

func indexHandler(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Title   string
			Content string
		}{
			Title:   "Welcome",
			Content: "<h1>Welcome to Can!</h1><p>A minimal web server with standard net/http.</p>",
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := IndexTemplate.Execute(w, data); err != nil {
			logger.Error("Failed to execute template", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}
