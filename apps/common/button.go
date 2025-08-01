package common

import "github.com/lumiknit/can/pkg/web"

type Button struct {
	Class string
	Label string
}

func (x *Button) Render(b *web.Builder) error {
	return b.P(
		`<button class="`, x.Class, `">`,
		web.Esc(x.Label),
		`</button>`,
	)
}
