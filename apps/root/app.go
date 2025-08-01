package root

import (
	"context"
	"net/http"

	"github.com/lumiknit/can/pkg/web"
)

func App() *web.App {
	app := &web.App{
		BasePath: "/",
	}

	// Add root page
	app.AddPage(&web.Page{
		Path:  "/",
		Title: "Welcome to Can!",
		Handle: func(ctx context.Context, req *http.Request) (web.Component, error) {
			return web.Tag("div", nil,
				web.Tag("h1", nil, "Welcome to Can!"),
				web.Tag("p", nil, "A minimal web server with standard net/http."),
				web.Tag("p", nil,
					"Explore our apps:",
				),
				web.Tag("ul", nil,
					web.Tag("li", nil, web.Tag("a", web.Attrs{"href": "/hello"}, "Hello World example")),
					web.Tag("li", nil, web.Tag("a", web.Attrs{"href": "/thread"}, "Thread Board")),
					web.Tag("li", nil, web.Tag("a", web.Attrs{"href": "/chatai"}, "Chat AI")),
				),
			), nil
		},
	})

	return app
}
