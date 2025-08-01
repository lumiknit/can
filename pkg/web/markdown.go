package web

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// MarkdownComponent renders markdown content as HTML
type MarkdownComponent struct {
	Content string
}

func Markdown(content string) *MarkdownComponent {
	return &MarkdownComponent{
		Content: content,
	}
}

func (m *MarkdownComponent) Render(b *Builder) error {
	// Configure markdown parser with common extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)

	// Configure HTML renderer
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	// Convert markdown to HTML
	htmlBytes := markdown.ToHTML([]byte(m.Content), p, renderer)

	// Output raw HTML (no escaping)
	return b.P(string(htmlBytes))
}
