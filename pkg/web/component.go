package web

import "strings"

// Component is a unit of web content, which can be rendered into HTML.
type Component interface {
	// Render renders the component into HTML.
	// It should use builder to append raw HTML content.
	Render(b *Builder) error
}

type Attrs map[string]string

// Tag represents a generic HTML tag.
type TagComponent struct {
	Name   string
	IsVoid bool
	Attrs
	Body []any
}

var _ Component = (*TagComponent)(nil)

// voidElements are HTML elements that cannot have content and must be self-closing
var voidElements = map[string]bool{
	"area":   true,
	"base":   true,
	"br":     true,
	"col":    true,
	"embed":  true,
	"hr":     true,
	"img":    true,
	"input":  true,
	"link":   true,
	"meta":   true,
	"source": true,
	"track":  true,
	"wbr":    true,
}

func Tag(name string, attrs Attrs, body ...any) *TagComponent {
	// Convert tag name to lowercase
	normalizedName := strings.ToLower(name)

	return &TagComponent{
		Name:   normalizedName,
		IsVoid: voidElements[normalizedName],
		Attrs:  attrs,
		Body:   body,
	}
}

func (t *TagComponent) Render(b *Builder) error {
	if err := b.P("<", t.Name); err != nil {
		return err
	}

	for k, v := range t.Attrs {
		if err := b.P(" ", k, `="`, Esc(v), `"`); err != nil {
			return err
		}
	}

	if t.IsVoid {
		return b.P(" />")
	}

	if err := b.P(">"); err != nil {
		return err
	}

	for _, body := range t.Body {
		if err := b.P(body); err != nil {
			return err
		}
	}

	return b.P("</", t.Name, ">")
}
