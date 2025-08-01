package hello

import (
	"github.com/lumiknit/can/pkg/web"
)

type pageTop struct{}

func (p *pageTop) Render(b *web.Builder) error {

	tag := web.Tag(
		"h1",
		web.Attrs{
			"class": "page-title",
		},
		web.Esc("Hello, world!"),
	)

	return b.P(
		tag,
	)
}
