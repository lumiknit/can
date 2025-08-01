package web

import (
	"context"
	"fmt"
)

func StaticPath(rest string) string {
	return "/_/" + rest
}

func APIPath(rest string) string {
	return "/-/" + rest
}

type RenderContext struct {
	URL      string
	BasePath string
}

// Renderable is an interface for types that can be rendered to HTML.
type Renderable interface {
	Render(ctx context.Context, rc RenderContext) ([]byte, error)
}

// Stylesheet represents a link to a CSS Stylesheet in a web page.
type Stylesheet struct {
	// Href is the path of the stylesheet.
	// If not present, the stylesheet will be inlined in the HTML.
	Href string

	// Content is the CSS content of the stylesheet.
	Content []byte
}

func (s *Stylesheet) Render(ctx context.Context, rc RenderContext) ([]byte, error) {
	if s.Href == "" {
		return fmt.Appendf(nil, "<style>%s</style>", s.Content), nil
	}
	// Otherwise, link to the external stylesheet.
	return fmt.Appendf(nil, `<link rel="stylesheet" href="%s">`, rc.BasePath+s.Href), nil
}
