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

	// Add common presets (favicon, icons, etc.)
	builder.AddCommonPresets()

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

	// HTML declaration with lang from builder
	lang := b.Lang
	if lang == "" {
		lang = "en" // fallback
	}
	html.WriteString(fmt.Sprintf("<!DOCTYPE html>\n<html lang=\"%s\"><head>", lang))

	// Meta tags from builder
	for _, meta := range b.MetaTags {
		meta.WriteHTML(&html)
	}

	// Link tags from builder
	for _, link := range b.LinkTags {
		link.WriteHTML(&html)
	}

	// Title
	if b.Title != "" {
		html.WriteString("<title>")
		html.WriteString(b.Title)
		html.WriteString("</title>")
	}

	// Stylesheet links
	for _, href := range b.Stylesheets {
		(&LinkTag{
			Rel:  "stylesheet",
			Href: href,
		}).WriteHTML(&html)
	}

	html.WriteString("</head><body>")

	// Body content
	html.WriteString(b.String())

	// Script files
	for _, src := range b.ScriptFiles {
		html.WriteString(fmt.Sprintf("  <script src=\"%s\"></script>", src))
	}

	// Inline scripts
	for _, script := range b.InlineScripts {
		html.WriteString("  <script>")
		html.WriteString(script)
		html.WriteString("</script>")
	}

	html.WriteString("</body></html>")

	return html.String()
}
