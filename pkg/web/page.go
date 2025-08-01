package web

import (
	"context"
	"net/http"
)

type PageHandler func(ctx context.Context, req *http.Request) (Component, error)

// Page represents a web page in the application.
// A single app contains multiple pages.
type Page struct {
	Parent *App
	Path   string

	Title       string
	Stylesheets []string // File paths to stylesheets

	Handle PageHandler
}
