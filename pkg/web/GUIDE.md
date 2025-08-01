# Web Framework Guide

The `web` package provides a framework for building static web applications in Go using the App-Page-Component architecture.

## Architecture Overview

### App → Page → Component

- **App**: The root of a web application, mounted to a specific path
- **Page**: Individual HTML pages within an app
- **Component**: Reusable HTML building blocks that render content

## Core Concepts

### App

An App represents a complete web application unit that can be mounted to a specific path in the HTTP server.

**Features:**
- Contains static files (CSS, JS, images, etc.)
- Has global stylesheets that apply to all pages
- Contains multiple pages
- Mounted to a base path (e.g., `/`, `/admin`, `/blog`)

**Example:**
```go
func NewApp() *web.App {
    app := &web.App{
        BasePath: "/myapp",
    }
    
    // Add static files
    app.AddStaticEmbedFS(staticFS, "static")
    
    // Add global stylesheets
    app.AddStyle(&web.Stylesheet{
        Href: "/static/global.css",
    })
    
    // Add pages
    app.AddPage(&web.Page{
        Path:  "/",
        Title: "Home",
        Handle: homeHandler,
    })
    
    return app
}
```

### Page

A Page represents a single HTML document with its own URL path, title, and content.

**Features:**
- Has a unique path within the app
- Can have page-specific stylesheets
- Contains a handler function that returns a Component
- Automatically gets wrapped in HTML5 document structure

**Example:**
```go
&web.Page{
    Path:  "/about",
    Title: "About Us",
    Stylesheets: []string{"/static/about.css"},
    Handle: func(ctx context.Context, req *http.Request) (web.Component, error) {
        return &aboutComponent{}, nil
    },
}
```

### Component

A Component is a reusable HTML building block that implements the `Component` interface.

**Interface:**
```go
type Component interface {
    Render(b *Builder) error
}
```

**Example:**
```go
type heroComponent struct {
    Title string
    Message string
}

func (h *heroComponent) Render(b *web.Builder) error {
    return b.P(
        web.Tag("div", web.Attrs{"class": "hero"},
            web.Tag("h1", nil, web.Esc(h.Title)),
            web.Tag("p", nil, web.Esc(h.Message)),
        ),
    )
}
```

## Project Structure

### Recommended App Structure

```
apps/
├── myapp/                 # App directory
│   ├── app.go            # App definition (NewApp function)
│   ├── p_home.go         # Page handlers (p_*.go)
│   ├── p_about.go        
│   ├── c_hero.go         # Components (c_*.go)
│   ├── c_footer.go       
│   └── static/           # Static files
│       ├── style.css
│       └── script.js
```

### Naming Conventions

- **App function**: `NewApp()` in `app.go`
- **Pages**: `p_*.go` (e.g., `p_home.go`, `p_contact.go`)
- **Components**: `c_*.go` (e.g., `c_navbar.go`, `c_button.go`)
- **Private implementation**: Use lowercase names for internal components/pages

### Example App Implementation

**apps/blog/app.go:**
```go
package blog

import (
    "embed"
    "github.com/lumiknit/can/pkg/web"
)

//go:embed static/*
var staticFS embed.FS

func NewApp() *web.App {
    app := &web.App{
        BasePath: "/blog",
    }
    
    // Add static files
    app.AddStaticEmbedFS(staticFS, "static")
    
    // Add global stylesheet
    app.AddStyle(&web.Stylesheet{
        Href: "/static/blog.css",
    })
    
    // Add pages
    app.AddPage(homePage())
    app.AddPage(postPage())
    
    return app
}
```

**apps/blog/p_home.go:**
```go
package blog

import (
    "context"
    "net/http"
    "github.com/lumiknit/can/pkg/web"
)

func homePage() *web.Page {
    return &web.Page{
        Path:  "/",
        Title: "My Blog",
        Handle: homeHandler,
    }
}

func homeHandler(ctx context.Context, req *http.Request) (web.Component, error) {
    return &homeComponent{
        Posts: []Post{
            {Title: "First Post", Content: "Hello world!"},
        },
    }, nil
}
```

**apps/blog/c_home.go:**
```go
package blog

import "github.com/lumiknit/can/pkg/web"

type homeComponent struct {
    Posts []Post
}

func (h *homeComponent) Render(b *web.Builder) error {
    return b.P(
        web.Tag("div", web.Attrs{"class": "home"},
            web.Tag("h1", nil, "Welcome to My Blog"),
            h.renderPosts(),
        ),
    )
}

func (h *homeComponent) renderPosts() web.Component {
    var posts []any
    for _, post := range h.Posts {
        posts = append(posts, &postPreview{Post: post})
    }
    return web.Tag("div", web.Attrs{"class": "posts"}, posts...)
}
```

## HTML Generation

### Builder API

The Builder provides methods for constructing HTML:

- `b.P(content...)` - Add content to the buffer
- `b.SetTitle(title)` - Set page title
- `b.AddStylesheet(href)` - Add stylesheet link
- `b.AddScriptFile(src)` - Add script file
- `b.AddInlineScript(script)` - Add inline script

### HTML Tags

Use `web.Tag()` to create HTML elements:

```go
web.Tag("div", web.Attrs{"class": "container", "id": "main"},
    web.Tag("h1", nil, "Hello World"),
    web.Tag("p", nil, "This is a paragraph."),
)
```

### Escaping

Use `web.Esc()` to escape HTML content:

```go
web.Tag("p", nil, web.Esc(userInput))
```

## Mounting Apps

Apps are mounted in the server router:

```go
// internal/server/router.go
func SetupRoutes(logger *slog.Logger, releaseMode bool) http.Handler {
    mux := http.NewServeMux()
    
    // Mount apps
    mountApp(mux, root.NewApp())
    mountApp(mux, blog.NewApp())
    
    return handler
}
```

## Error Handling

Use predefined web errors:

```go
import "github.com/lumiknit/can/pkg/web"

func handler(ctx context.Context, req *http.Request) (web.Component, error) {
    if unauthorized {
        return nil, web.ErrUnauthorized
    }
    
    if notFound {
        return nil, web.ErrNotFound
    }
    
    // Custom error
    return nil, web.NewWebError(422, "Invalid input")
}
```

## Best Practices

1. **Keep components small and focused** - Each component should have a single responsibility
2. **Use private types** - Keep internal components and pages private (lowercase names)
3. **Embed static files** - Use `go:embed` for CSS, JS, and other assets
4. **Handle errors gracefully** - Return appropriate web errors from handlers
5. **Use semantic HTML** - Build accessible, standards-compliant markup
6. **Organize by feature** - Group related pages and components together
7. **Test components** - Components are easy to unit test in isolation

## Example: Complete Minimal App

```go
// apps/example/app.go
package example

import (
    "context"
    "net/http"
    "github.com/lumiknit/can/pkg/web"
)

func NewApp() *web.App {
    app := &web.App{BasePath: "/example"}
    
    app.AddPage(&web.Page{
        Path:  "/",
        Title: "Example App",
        Handle: func(ctx context.Context, req *http.Request) (web.Component, error) {
            return web.Tag("div", nil,
                web.Tag("h1", nil, "Hello from Example App!"),
                web.Tag("p", nil, "This is a minimal example."),
            ), nil
        },
    })
    
    return app
}
```

This creates a complete app accessible at `/example/` with a simple HTML page.