package server

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lumiknit/can/pkg/css"
	"github.com/lumiknit/can/static"
)

func SetupRoutes(logger *slog.Logger, releaseMode bool) *gin.Engine {
	if releaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	static.Each(func(file *static.Item) {
		data := file.Content
		switch strings.Split(file.ContentType, ";")[0] {
		case "text/css":
			if stylesheet, err := css.Parse(string(file.Content)); err == nil {
				data = []byte(stylesheet.Minify())
			} else {
				logger.Warn("Failed to parse CSS for minification", "path", file.Path, "error", err)
			}
		}

		r.GET(file.Path, func(c *gin.Context) {
			c.Header("Content-Type", file.ContentType)
			c.Data(http.StatusOK, file.ContentType, data)
		})
		r.HEAD(file.Path, func(c *gin.Context) {
			c.Header("Content-Type", file.ContentType)
			c.Status(http.StatusOK)
		})
	})

	// Health check endpoint
	r.GET("/healthz", healthzHandler)

	// Index page
	r.GET("/", indexHandler(logger))

	return r
}

func indexHandler(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		data := struct {
			Title   string
			Content string
		}{
			Title:   "Welcome",
			Content: "<h1>Welcome to Can!</h1><p>A minimal web server with Gin framework.</p>",
		}

		c.Header("Content-Type", "text/html; charset=utf-8")
		if err := IndexTemplate.Execute(c.Writer, data); err != nil {
			logger.Error("Failed to execute template", "error", err)
			c.String(http.StatusInternalServerError, "Internal Server Error")
			return
		}
	}
}
