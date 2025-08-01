package hello

import (
	"embed"

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
	a := &web.App{}
	a.AddStaticEmbedFS(fs, "static")
	a.AddStyle(MyCSS())

	return a
}
