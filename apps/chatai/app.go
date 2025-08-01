package chatai

import (
	"time"

	"github.com/lumiknit/can/pkg/web"
)

// Message represents a chat message
type Message struct {
	ID        int       `json:"id"`
	Role      string    `json:"role"` // "user" or "assistant"
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// ChatStore manages chat messages in memory
type ChatStore struct {
	messages     []Message
	nextID       int
	client       *OpenAIClient
	model        string
	isProcessing bool
}

var store = &ChatStore{
	messages:     []Message{},
	nextID:       1,
	client:       NewOpenAIClient(),
	model:        "gpt-4o-mini", // Default model
	isProcessing: false,
}

// AddUserMessage adds a user message and starts background AI processing
func (s *ChatStore) AddUserMessage(content string) error {
	// Add user message
	userMsg := Message{
		ID:        s.nextID,
		Role:      "user",
		Content:   content,
		CreatedAt: time.Now(),
	}
	s.messages = append(s.messages, userMsg)
	s.nextID++

	// Start background AI processing
	s.isProcessing = true
	go s.processAIResponse()

	return nil
}

// processAIResponse handles AI response in background
func (s *ChatStore) processAIResponse() {
	defer func() {
		s.isProcessing = false
	}()

	// Get AI response if client is available
	if s.client != nil {
		// Convert to OpenAI format
		var chatMessages []ChatMessage
		for _, msg := range s.messages {
			chatMessages = append(chatMessages, ChatMessage{
				Role:    msg.Role,
				Content: msg.Content,
			})
		}

		resp, err := s.client.ChatCompletion(chatMessages, s.model)
		if err != nil {
			// Add error message
			errorMsg := Message{
				ID:        s.nextID,
				Role:      "assistant",
				Content:   "Error: " + err.Error(),
				CreatedAt: time.Now(),
			}
			s.messages = append(s.messages, errorMsg)
			s.nextID++
			return
		}

		if len(resp.Choices) > 0 {
			// Add AI response
			aiMsg := Message{
				ID:        s.nextID,
				Role:      "assistant",
				Content:   resp.Choices[0].Message.Content,
				CreatedAt: time.Now(),
			}
			s.messages = append(s.messages, aiMsg)
			s.nextID++
		}
	} else {
		// Add fallback message if no API key
		fallbackMsg := Message{
			ID:        s.nextID,
			Role:      "assistant",
			Content:   "OpenAI API key not configured. Set OPENAI_API_KEY environment variable.",
			CreatedAt: time.Now(),
		}
		s.messages = append(s.messages, fallbackMsg)
		s.nextID++
	}
}

// GetMessages returns all messages
func (s *ChatStore) GetMessages() []Message {
	return s.messages
}

// GetModel returns current model
func (s *ChatStore) GetModel() string {
	return s.model
}

// SetModel sets the model
func (s *ChatStore) SetModel(model string) {
	s.model = model
}

// IsProcessing returns true if AI is currently processing a response
func (s *ChatStore) IsProcessing() bool {
	return s.isProcessing
}

// ClearMessages clears all messages
func (s *ChatStore) ClearMessages() {
	s.messages = []Message{}
	s.nextID = 1
}

func NewApp() *web.App {
	app := &web.App{
		BasePath: "/chatai",
	}

	// Add chat page
	app.AddPage(chatPage())

	return app
}
