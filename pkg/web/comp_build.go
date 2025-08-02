package web

import (
	"context"
	"fmt"
	"html"
	"io"
	"strconv"
	"strings"
)

// Esc is a string type, which represents escaped HTML content.
// When used in the Builder, it will be automatically escaped to prevent XSS attacks.
type Esc string

// Builder is a structure that helps in building text contents for HTML rendering.
type Builder struct {
	// Ctx is the current request context.
	Ctx context.Context

	// Buf is the buffer that accumulates the HTML content.
	Buf strings.Builder

	// Title for the HTML document
	Title string

	// Lang for the HTML document
	Lang string

	// Meta tags to include in the head
	MetaTags []MetaTag

	// Link tags to include in the head (rel, href, etc.)
	LinkTags []LinkTag

	// Stylesheets to include as <link> tags
	Stylesheets []string

	// Inline scripts to include in the body
	InlineScripts []string

	// Script files to include as <script src> tags
	ScriptFiles []string
}

// MetaTag represents an HTML meta tag with various attributes.
type MetaTag struct {
	Name      string
	Content   string
	HTTPEquiv string
	Property  string
	Charset   string
	Media     string
}

func (m *MetaTag) WriteHTML(w io.Writer) {
	w.Write([]byte("<meta"))
	if m.Name != "" {
		w.Write(fmt.Appendf(nil, " name=\"%s\"", html.EscapeString(m.Name)))
	}
	if m.Content != "" {
		w.Write(fmt.Appendf(nil, " content=\"%s\"", html.EscapeString(m.Content)))
	}
	if m.HTTPEquiv != "" {
		w.Write(fmt.Appendf(nil, " http-equiv=\"%s\"", html.EscapeString(m.HTTPEquiv)))
	}
	if m.Property != "" {
		w.Write(fmt.Appendf(nil, " property=\"%s\"", html.EscapeString(m.Property)))
	}
	if m.Charset != "" {
		w.Write(fmt.Appendf(nil, " charset=\"%s\"", html.EscapeString(m.Charset)))
	}
	if m.Media != "" {
		w.Write(fmt.Appendf(nil, " media=\"%s\"", html.EscapeString(m.Media)))
	}
	w.Write([]byte(">"))
}

// LinkTag represents an HTML link tag with various attributes.
type LinkTag struct {
	Rel   string
	Href  string
	Type  string
	Sizes string
	Media string
}

func (l *LinkTag) WriteHTML(w io.Writer) {
	w.Write([]byte("<link"))
	if l.Rel != "" {
		w.Write(fmt.Appendf(nil, " rel=\"%s\"", html.EscapeString(l.Rel)))
	}
	if l.Href != "" {
		w.Write(fmt.Appendf(nil, " href=\"%s\"", html.EscapeString(l.Href)))
	}
	if l.Type != "" {
		w.Write(fmt.Appendf(nil, " type=\"%s\"", html.EscapeString(l.Type)))
	}
	if l.Sizes != "" {
		w.Write(fmt.Appendf(nil, " sizes=\"%s\"", html.EscapeString(l.Sizes)))
	}
	if l.Media != "" {
		w.Write(fmt.Appendf(nil, " media=\"%s\"", html.EscapeString(l.Media)))
	}
	w.Write([]byte(">"))
}

func NewBuilder(ctx context.Context) *Builder {
	return &Builder{
		Ctx:  ctx,
		Buf:  strings.Builder{},
		Lang: "en", // Default language
		MetaTags: []MetaTag{
			{Name: "viewport", Content: "width=device-width, initial-scale=1.0"},
			{Name: "charset", Content: "UTF-8"},
		},
		LinkTags:      []LinkTag{},
		Stylesheets:   []string{},
		InlineScripts: []string{},
		ScriptFiles:   []string{},
	}
}

// String returns the accumulated HTML content as a string.
func (b *Builder) String() string {
	return b.Buf.String()
}

// AddStylesheet adds a stylesheet link to the builder.
func (b *Builder) AddStylesheet(href string) {
	if href != "" {
		b.Stylesheets = append(b.Stylesheets, href)
	}
}

// AddInlineScript adds an inline script to the builder.
// This script will be included in the HTML body.
func (b *Builder) AddInlineScript(script string) {
	if script != "" {
		b.InlineScripts = append(b.InlineScripts, script)
	}
}

// AddScriptFile adds a script file to the builder.
// This will be included as a <script src> tag in the HTML.
func (b *Builder) AddScriptFile(src string) {
	if src != "" {
		b.ScriptFiles = append(b.ScriptFiles, src)
	}
}

// SetTitle sets the title of the HTML document.
// This will be used in the <title> tag of the HTML head.
func (b *Builder) SetTitle(title string) {
	b.Title = title
}

// SetLang sets the language of the HTML document.
func (b *Builder) SetLang(lang string) {
	b.Lang = lang
}

// AddMetaTag adds a meta tag to the builder.
func (b *Builder) AddMetaTag(meta MetaTag) {
	b.MetaTags = append(b.MetaTags, meta)
}

// AddLinkTag adds a link tag to the builder.
func (b *Builder) AddLinkTag(link LinkTag) {
	b.LinkTags = append(b.LinkTags, link)
}

// SetRefreshMeta sets a meta tag for HTTP refresh.
func (b *Builder) SetRefreshMeta(seconds int, url string) {
	content := strconv.Itoa(seconds)
	if url != "" {
		content += "; url=" + url
	}
	b.AddMetaTag(MetaTag{
		HTTPEquiv: "refresh",
		Content:   content,
	})
}

// SetCharsetMeta sets a meta tag for character set.
func (b *Builder) SetCharsetMeta(charset string) {
	if charset != "" {
		b.AddMetaTag(MetaTag{
			Charset: charset,
		})
	}
}

// SetContentTypeMeta sets a meta tag for content type.
func (b *Builder) SetContentTypeMeta(contentType string) {
	if contentType != "" {
		b.AddMetaTag(MetaTag{
			Name:    "Content-Type",
			Content: contentType,
		})
	}
}

// SetThemeColorMeta sets a meta tag for theme color.
func (b *Builder) SetThemeColorMeta(color string, darkModeColor string) {
	if darkModeColor != "" {
		b.AddMetaTag(MetaTag{
			Name:    "theme-color",
			Content: darkModeColor,
			Media:   "(prefers-color-scheme: dark)",
		})
	}
	if color != "" {
		b.AddMetaTag(MetaTag{
			Name:    "theme-color",
			Content: color,
		})
	}
}

// AddCommonPresets adds common favicon and icon links.
func (b *Builder) AddCommonPresets() {
	// Favicon and icon links
	b.AddLinkTag(LinkTag{Rel: "icon", Type: "image/png", Href: "/favicon-96x96.png", Sizes: "96x96"})
	b.AddLinkTag(LinkTag{Rel: "icon", Type: "image/svg+xml", Href: "/favicon.svg"})
	b.AddLinkTag(LinkTag{Rel: "shortcut icon", Href: "/favicon.ico"})
	b.AddLinkTag(LinkTag{Rel: "apple-touch-icon", Href: "/apple-touch-icon.png", Sizes: "180x180"})
	b.AddLinkTag(LinkTag{Rel: "manifest", Href: "/manifest.json"})

	// Common meta tags
	b.AddMetaTag(MetaTag{Name: "apple-mobile-web-app-title", Content: "Can"})
	b.SetThemeColorMeta("#000000", "")
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
