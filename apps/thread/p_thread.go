package thread

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/lumiknit/can/pkg/web"
)

func threadPage() *web.Page {
	return &web.Page{
		Path:   "/",
		Title:  "Thread Board",
		Handle: threadHandler,
	}
}

func threadHandler(ctx context.Context, req *http.Request) (web.Component, error) {
	// Handle POST request for creating new posts
	if req.Method == http.MethodPost {
		if err := req.ParseForm(); err != nil {
			return nil, web.ErrBadRequest
		}

		content := req.FormValue("content")
		if strings.TrimSpace(content) == "" {
			return nil, web.NewWebError(400, "Content cannot be empty")
		}

		store.AddPost(content)

		// Redirect to prevent duplicate submissions
		// Since we can't redirect in component, we'll just show the page
	}

	// Parse pagination parameters
	skip := 0
	size := 10

	if skipStr := req.URL.Query().Get("skip"); skipStr != "" {
		if parsed, err := strconv.Atoi(skipStr); err == nil && parsed >= 0 {
			skip = parsed
		}
	}

	if sizeStr := req.URL.Query().Get("size"); sizeStr != "" {
		if parsed, err := strconv.Atoi(sizeStr); err == nil && parsed > 0 && parsed <= 100 {
			size = parsed
		}
	}

	// Get posts
	posts := store.GetPosts(skip, size)
	totalCount := store.GetTotalCount()

	return &threadComponent{
		Posts:      posts,
		Skip:       skip,
		Size:       size,
		TotalCount: totalCount,
	}, nil
}
