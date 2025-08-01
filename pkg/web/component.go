package web

// Component is a unit of web content, which can be rendered into HTML.
type Component interface {
	// Render renders the component into HTML.
	// It should use builder to append raw HTML content.
	Render(b *Builder) error
}

type Attrs map[string]string

// Tag represents a generic HTML tag.
type TagComponent struct {
	Name string
	Attrs
	Body []any
}

var _ Component = (*TagComponent)(nil)

func Tag(name string, attrs Attrs, body ...any) *TagComponent {
	return &TagComponent{
		Name:  name,
		Attrs: attrs,
		Body:  body,
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

	if len(t.Body) == 0 {
		return b.P("/>")
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
