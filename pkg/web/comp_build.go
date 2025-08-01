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

	// Title for the HTML document
	Title string

	// Meta tags to include in the head
	MetaTags []MetaTag

	// Stylesheets to include as <link> tags
	Stylesheets []string

	// Inline scripts to include in the body
	InlineScripts []string

	// Script files to include as <script src> tags
	ScriptFiles []string
}

type MetaTag struct {
	Name      string
	Content   string
	HTTPEquiv string
	Property  string
	Charset   string
}

func NewBuilder(ctx context.Context) *Builder {
	return &Builder{
		Ctx:           ctx,
		Buf:           strings.Builder{},
		MetaTags:      []MetaTag{},
		Stylesheets:   []string{},
		InlineScripts: []string{},
		ScriptFiles:   []string{},
	}
}

func (b *Builder) String() string {
	return b.Buf.String()
}

func (b *Builder) AddStylesheet(href string) {
	if href != "" {
		b.Stylesheets = append(b.Stylesheets, href)
	}
}

func (b *Builder) AddInlineScript(script string) {
	if script != "" {
		b.InlineScripts = append(b.InlineScripts, script)
	}
}

func (b *Builder) AddScriptFile(src string) {
	if src != "" {
		b.ScriptFiles = append(b.ScriptFiles, src)
	}
}

func (b *Builder) SetTitle(title string) {
	b.Title = title
}

func (b *Builder) AddMetaTag(meta MetaTag) {
	b.MetaTags = append(b.MetaTags, meta)
}

func (b *Builder) SetRefresh(seconds int, url string) {
	content := strconv.Itoa(seconds)
	if url != "" {
		content += "; url=" + url
	}
	b.AddMetaTag(MetaTag{
		HTTPEquiv: "refresh",
		Content:   content,
	})
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
