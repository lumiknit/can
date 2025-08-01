package thread

import (
	"fmt"
	"strconv"

	"github.com/lumiknit/can/pkg/web"
)

type threadComponent struct {
	Posts      []Post
	Skip       int
	Size       int
	TotalCount int
}

func (t *threadComponent) Render(b *web.Builder) error {
	return b.P(
		web.Tag("div", nil,
			web.Tag("h1", nil, "Thread Board"),

			// Post form
			web.Tag("form", web.Attrs{"method": "post", "action": ""},
				web.Tag("div", nil,
					web.Tag("textarea", web.Attrs{
						"name":        "content",
						"placeholder": "Write your message here...",
						"rows":        "4",
						"cols":        "80",
						"required":    "required",
					}),
				),
				web.Tag("div", nil,
					web.Tag("button", web.Attrs{"type": "submit"}, "Post"),
				),
			),

			web.Tag("hr", nil),

			// Posts list
			t.renderPosts(),

			// Pagination
			t.renderPagination(),
		),
	)
}

func (t *threadComponent) renderPosts() web.Component {
	if len(t.Posts) == 0 {
		return web.Tag("p", nil, "No posts yet. Be the first to post!")
	}

	var posts []any
	for _, post := range t.Posts {
		posts = append(posts, &postComponent{Post: post})
	}

	return web.Tag("div", nil, posts...)
}

func (t *threadComponent) renderPagination() web.Component {
	if t.TotalCount <= t.Size {
		return web.Tag("div", nil) // Empty div if no pagination needed
	}

	var links []any

	// Previous link
	if t.Skip > 0 {
		prevSkip := t.Skip - t.Size
		if prevSkip < 0 {
			prevSkip = 0
		}
		prevURL := fmt.Sprintf("?skip=%d&size=%d", prevSkip, t.Size)
		links = append(links,
			web.Tag("a", web.Attrs{"href": prevURL}, "← Previous"),
			Text(" "),
		)
	}

	// Page info
	currentPage := (t.Skip / t.Size) + 1
	totalPages := (t.TotalCount + t.Size - 1) / t.Size
	links = append(links,
		Text(fmt.Sprintf("Page %d of %d", currentPage, totalPages)),
		Text(" "),
	)

	// Next link
	if t.Skip+t.Size < t.TotalCount {
		nextSkip := t.Skip + t.Size
		nextURL := fmt.Sprintf("?skip=%d&size=%d", nextSkip, t.Size)
		links = append(links,
			web.Tag("a", web.Attrs{"href": nextURL}, "Next →"),
		)
	}

	return web.Tag("div", nil, links...)
}

type postComponent struct {
	Post Post
}

func (p *postComponent) Render(b *web.Builder) error {
	timeStr := p.Post.CreatedAt.Format("2006-01-02 15:04:05")

	return b.P(
		web.Tag("div", web.Attrs{"style": "border: 1px solid #ccc; margin: 10px 0; padding: 10px; border-radius: 5px;"},
			web.Tag("div", web.Attrs{"style": "font-size: 0.9em; color: #666; margin-bottom: 5px;"},
				Text("#"), Text(strconv.Itoa(p.Post.ID)),
				Text(" - "), Text(timeStr),
			),
			web.Tag("div", web.Attrs{"style": "line-height: 1.6;"},
				web.Markdown(p.Post.Content),
			),
		),
	)
}

// textNode creates a simple text component
type textNode struct {
	text string
}

func (t *textNode) Render(b *web.Builder) error {
	return b.P(t.text)
}

func Text(s string) web.Component {
	return &textNode{text: s}
}
