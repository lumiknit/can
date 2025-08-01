package web

import "embed"

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
