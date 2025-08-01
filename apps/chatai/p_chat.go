package chatai

import (
	"context"
	"net/http"
	"strings"

	"github.com/lumiknit/can/pkg/web"
)

func chatPage() *web.Page {
	return &web.Page{
		Path:   "/",
		Title:  "Chat AI",
		Handle: chatHandler,
	}
}

func chatHandler(ctx context.Context, req *http.Request) (web.Component, error) {
	// Handle POST request for sending messages
	if req.Method == http.MethodPost {
		if err := req.ParseForm(); err != nil {
			return nil, web.ErrBadRequest
		}

		action := req.FormValue("action")

		switch action {
		case "send":
			message := strings.TrimSpace(req.FormValue("message"))
			if message == "" {
				return nil, web.NewWebError(400, "Message cannot be empty")
			}

			if err := store.AddUserMessage(message); err != nil {
				// Error is already added to messages, continue to show page
			}
		case "clear":
			store.ClearMessages()
		case "model":
			model := req.FormValue("model")
			if model != "" {
				store.SetModel(model)
			}
		}
	}

	// Get messages and processing state
	messages := store.GetMessages()
	isProcessing := store.IsProcessing()

	return &chatComponent{
		Messages:     messages,
		CurrentModel: store.GetModel(),
		IsProcessing: isProcessing,
	}, nil
}
