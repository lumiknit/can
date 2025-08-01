package web

import (
	"context"
	"html"
	"strconv"
	"strings"
)

type Esc string

// Builder is a structure that helps in building text contents for HTML rendering.

type Builder struct {
	Ctx context.Context
	Buf strings.Builder

	Stylesheets []*Stylesheet
}

func NewBuilder(ctx context.Context) *Builder {
	return &Builder{
		Ctx: ctx,
		Buf: strings.Builder{},
	}
}

func (b *Builder) String() string {
	return b.Buf.String()
}

func (b *Builder) AddStylesheet(s *Stylesheet) {
	if s != nil {
		b.Stylesheets = append(b.Stylesheets, s)
	}
}

// P pushes a contents to the builder's buffer.
func (b *Builder) P(
	contents ...any,
) error {
	// Check if the context is done before proceeding
	if err := b.Ctx.Err(); err != nil {
		return err
	}

	for _, c := range contents {
		switch v := c.(type) {
		case string:
			b.Buf.WriteString(v)
		case []byte:
			b.Buf.Write(v)
		case Esc:
			b.Buf.WriteString(html.EscapeString(string(v)))
		case Component:
			if err := v.Render(b); err != nil {
				return err
			}
		case int:
			b.Buf.WriteString(strconv.Itoa(v))
		case float32:
		case float64:
			b.Buf.WriteString(strconv.FormatFloat(float64(v), 'f', -1, 32))
		case bool:
			if v {
				b.Buf.WriteString("true")
			} else {
				b.Buf.WriteString("false")
			}
		case []any:
			b.P(v...)
		}
	}

	return nil
}
