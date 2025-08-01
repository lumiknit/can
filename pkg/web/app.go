package web

import (
	"embed"
	"fmt"
	"net/http"
	"strings"
)

// App is used to define a web application unit.
// Each app may contain files, stylesheets and pages.
// App can be mounted to a specific path in the http server.
type App struct {
	BasePath string

	Manifest *Manifest

	// StaticFiles is a set of static files which can be served by the application.
	StaticFiles []*StaticFile

	// Stylesheets is a set of stylesheets which applies over the entire application.
	Stylesheets []*Stylesheet

	// Pages is a set of web pages which can be rendered.
	Pages []*Page
}

// AddStyle adds a stylesheet to the application.
func (a *App) AddStyle(s *Stylesheet) {
	a.Stylesheets = append(a.Stylesheets, s)
}

// AddStaticEmbedFS adds a static file system to the application.
func (a *App) AddStaticEmbedFS(efs embed.FS, rootDir string) error {
	files, err := EmbedFSToStaticFiles(efs, rootDir)
	if err != nil {
		return err
	}
	a.StaticFiles = append(a.StaticFiles, files...)
	return nil
}

func (a *App) AddPage(p *Page) {
	p.Parent = a
	a.Pages = append(a.Pages, p)
}

// HandlePage handles HTTP requests by finding the appropriate page and rendering it
func (a *App) HandlePage(w http.ResponseWriter, r *http.Request) {
	// Find the matching page
	// Remove the app's base path from the request URL
	requestPath := strings.TrimPrefix(r.URL.Path, a.BasePath)
	if requestPath == "" {
		requestPath = "/"
	}

	var matchedPage *Page
	for _, page := range a.Pages {
		// Normalize paths for comparison (remove trailing slash except for root)
		pagePath := page.Path
		if pagePath != "/" && strings.HasSuffix(pagePath, "/") {
			pagePath = strings.TrimSuffix(pagePath, "/")
		}

		normalizedRequestPath := requestPath
		if normalizedRequestPath != "/" && strings.HasSuffix(normalizedRequestPath, "/") {
			normalizedRequestPath = strings.TrimSuffix(normalizedRequestPath, "/")
		}

		if pagePath == normalizedRequestPath {
			matchedPage = page
			break
		}
	}

	if matchedPage == nil {
		ErrNotFound.WriteError(w)
		return
	}

	// Create context
	ctx := r.Context()

	// Call the page handler to get the component
	component, err := matchedPage.Handle(ctx, r)
	if err != nil {
		if webErr, ok := err.(*WebError); ok {
			webErr.WriteError(w)
		} else {
			ErrInternalServerError.WriteError(w)
		}
		return
	}

	// Create builder and set up global stylesheets
	builder := NewBuilder(ctx)
	builder.SetTitle(matchedPage.Title)

	// Add global stylesheets from app
	for _, stylesheet := range a.Stylesheets {
		if stylesheet.Href != "" {
			builder.AddStylesheet(a.BasePath + stylesheet.Href)
		}
	}

	// Add page-specific stylesheets
	for _, stylesheetPath := range matchedPage.Stylesheets {
		builder.AddStylesheet(a.BasePath + stylesheetPath)
	}

	// Render the component
	if err := component.Render(builder); err != nil {
		ErrInternalServerError.WriteError(w)
		return
	}

	// Build the full HTML document
	html := a.buildHTML5Document(builder)

	// Write response
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

// buildHTML5Document creates a complete HTML5 document from the builder
func (a *App) buildHTML5Document(b *Builder) string {
	var html strings.Builder

	html.WriteString("<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n")
	html.WriteString("  <meta charset=\"UTF-8\">\n")
	html.WriteString("  <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n")

	// Additional meta tags
	for _, meta := range b.MetaTags {
		html.WriteString("  <meta")
		if meta.Name != "" {
			html.WriteString(fmt.Sprintf(" name=\"%s\"", meta.Name))
		}
		if meta.HTTPEquiv != "" {
			html.WriteString(fmt.Sprintf(" http-equiv=\"%s\"", meta.HTTPEquiv))
		}
		if meta.Property != "" {
			html.WriteString(fmt.Sprintf(" property=\"%s\"", meta.Property))
		}
		if meta.Charset != "" {
			html.WriteString(fmt.Sprintf(" charset=\"%s\"", meta.Charset))
		}
		if meta.Content != "" {
			html.WriteString(fmt.Sprintf(" content=\"%s\"", meta.Content))
		}
		html.WriteString(">\n")
	}

	// Title
	if b.Title != "" {
		html.WriteString(fmt.Sprintf("  <title>%s</title>\n", b.Title))
	}

	// Stylesheet links
	for _, href := range b.Stylesheets {
		html.WriteString(fmt.Sprintf("  <link rel=\"stylesheet\" href=\"%s\">\n", href))
	}

	html.WriteString("</head>\n<body>\n")

	// Body content
	html.WriteString(b.String())

	// Script files
	for _, src := range b.ScriptFiles {
		html.WriteString(fmt.Sprintf("  <script src=\"%s\"></script>\n", src))
	}

	// Inline scripts
	for _, script := range b.InlineScripts {
		html.WriteString("  <script>\n")
		html.WriteString(script)
		html.WriteString("\n  </script>\n")
	}

	html.WriteString("</body>\n</html>")

	return html.String()
}
