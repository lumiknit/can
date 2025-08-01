package hello

import (
	"context"
	"embed"
	"net/http"

	"github.com/lumiknit/can/pkg/web"
)

//go:embed static/*
var fs embed.FS

func MyCSS() *web.Stylesheet {
	return &web.Stylesheet{
		Href:    "/static/global.css",
		Content: []byte("body { font-family: Arial, sans-serif; }"),
	}
}

func App() *web.App {
	a := &web.App{
		BasePath: "/hello",
	}
	a.AddStaticEmbedFS(fs, "static")
	a.AddStyle(MyCSS())

	// Add hello page
	a.AddPage(&web.Page{
		Path:  "/",
		Title: "Hello World",
		Handle: func(ctx context.Context, req *http.Request) (web.Component, error) {
			return &pageTop{}, nil
		},
	})

	return a
}
