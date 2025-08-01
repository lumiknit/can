package thread

import (
	"strings"
	"time"

	"github.com/lumiknit/can/pkg/web"
)

// Post represents a thread post
type Post struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// ThreadStore manages posts in memory
type ThreadStore struct {
	posts  []Post
	nextID int
}

var store = &ThreadStore{
	posts:  []Post{},
	nextID: 1,
}

// AddPost adds a new post to the store
func (s *ThreadStore) AddPost(content string) Post {
	post := Post{
		ID:        s.nextID,
		Content:   strings.TrimSpace(content),
		CreatedAt: time.Now(),
	}
	s.posts = append([]Post{post}, s.posts...) // Add to beginning
	s.nextID++
	return post
}

// GetPosts returns posts with pagination
func (s *ThreadStore) GetPosts(skip, size int) []Post {
	if skip >= len(s.posts) {
		return []Post{}
	}

	end := skip + size
	if end > len(s.posts) {
		end = len(s.posts)
	}

	return s.posts[skip:end]
}

// GetTotalCount returns total number of posts
func (s *ThreadStore) GetTotalCount() int {
	return len(s.posts)
}

func NewApp() *web.App {
	app := &web.App{
		BasePath: "/thread",
	}

	// Add thread page
	app.AddPage(threadPage())

	return app
}
