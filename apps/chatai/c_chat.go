package chatai

import (
	"fmt"
	"strconv"

	"github.com/lumiknit/can/pkg/web"
)

type chatComponent struct {
	Messages     []Message
	CurrentModel string
	IsProcessing bool
}

func (c *chatComponent) Render(b *web.Builder) error {
	var lastMessageID int
	if len(c.Messages) > 0 {
		lastMessageID = c.Messages[len(c.Messages)-1].ID
	}

	// Set auto-refresh when processing
	if c.IsProcessing {
		b.SetRefreshMeta(5, "")
	}

	return b.P(
		web.Tag("div", nil,
			web.Tag("h1", nil, "Chat AI"),

			// Model selector
			c.renderModelSelector(),

			// Processing indicator
			c.renderProcessingIndicator(),

			// Chat history
			c.renderMessages(),

			// Message form
			web.Tag("form", web.Attrs{"method": "post", "action": ""},
				web.Tag("input", web.Attrs{"type": "hidden", "name": "action", "value": "send"}),
				web.Tag("div", web.Attrs{"style": "display: flex; gap: 10px; margin-top: 20px;"},
					web.Tag("input", web.Attrs{
						"type":        "text",
						"name":        "message",
						"placeholder": "Type your message here...",
						"style":       "flex: 1; padding: 10px; border: 1px solid #ccc; border-radius: 4px;",
						"required":    "required",
					}),
					web.Tag("button", web.Attrs{
						"type":  "submit",
						"style": "padding: 10px 20px; background: #007bff; color: white; border: none; border-radius: 4px; cursor: pointer;",
					}, "Send"),
				),
			),

			// Clear button
			web.Tag("form", web.Attrs{"method": "post", "action": "", "style": "margin-top: 10px;"},
				web.Tag("input", web.Attrs{"type": "hidden", "name": "action", "value": "clear"}),
				web.Tag("button", web.Attrs{
					"type":  "submit",
					"style": "padding: 5px 15px; background: #dc3545; color: white; border: none; border-radius: 4px; cursor: pointer;",
				}, "Clear Chat"),
			),

			// Auto-scroll script
			c.renderAutoScrollScript(lastMessageID),
		),
	)
}

func (c *chatComponent) renderModelSelector() web.Component {
	models := []struct {
		Value string
		Label string
	}{
		{"gpt-4o", "GPT-4o"},
		{"gpt-4o-mini", "GPT-4o Mini"},
		{"gpt-4", "GPT-4"},
		{"gpt-4-turbo", "GPT-4 Turbo"},
		{"gpt-3.5-turbo", "GPT-3.5 Turbo"},
	}

	var options []any
	for _, model := range models {
		attrs := web.Attrs{"value": model.Value}
		if model.Value == c.CurrentModel {
			attrs["selected"] = "selected"
		}
		options = append(options, web.Tag("option", attrs, model.Label))
	}

	return web.Tag("div", web.Attrs{"style": "margin: 10px 0; padding: 10px; background: #f8f9fa; border-radius: 4px;"},
		web.Tag("form", web.Attrs{"method": "post", "action": "", "style": "display: flex; align-items: center; gap: 10px;"},
			web.Tag("input", web.Attrs{"type": "hidden", "name": "action", "value": "model"}),
			web.Tag("label", web.Attrs{"for": "model"}, "Model: "),
			web.Tag("select", web.Attrs{
				"name":     "model",
				"id":       "model",
				"onchange": "this.form.submit()",
				"style":    "padding: 5px; border: 1px solid #ccc; border-radius: 4px;",
			}, options...),
			web.Tag("span", web.Attrs{"style": "color: #666; font-size: 0.9em;"},
				"Current: ", Text(c.CurrentModel)),
		),
	)
}

func (c *chatComponent) renderProcessingIndicator() web.Component {
	if !c.IsProcessing {
		return Text("")
	}

	return &processingIndicatorComponent{}
}

type processingIndicatorComponent struct{}

func (p *processingIndicatorComponent) Render(b *web.Builder) error {
	// Add CSS for spinner animation
	b.AddInlineScript(`
var style = document.createElement('style');
style.textContent = '@keyframes spin { 0% { transform: rotate(0deg); } 100% { transform: rotate(360deg); } }';
document.head.appendChild(style);
	`)

	return b.P(
		web.Tag("div", web.Attrs{
			"style": "margin: 10px 0; padding: 10px; background: #fff3cd; border: 1px solid #ffeaa7; border-radius: 4px; color: #856404; text-align: center;",
		},
			web.Tag("div", web.Attrs{"style": "display: flex; align-items: center; justify-content: center; gap: 10px;"},
				web.Tag("div", web.Attrs{
					"style": "width: 16px; height: 16px; border: 2px solid #856404; border-top: 2px solid transparent; border-radius: 50%; animation: spin 1s linear infinite;",
				}),
				Text("AI is processing your message... Page will refresh automatically."),
			),
		),
	)
}

func (c *chatComponent) renderMessages() web.Component {
	if len(c.Messages) == 0 {
		return web.Tag("div", web.Attrs{"style": "padding: 20px; text-align: center; color: #666;"},
			"No messages yet. Start a conversation!",
		)
	}

	var messages []any
	for _, msg := range c.Messages {
		messages = append(messages, &messageComponent{Message: msg})
	}

	return web.Tag("div", web.Attrs{
		"style": "max-height: 400px; overflow-y: auto; border: 1px solid #ccc; padding: 10px; margin: 20px 0; border-radius: 4px; background: #f9f9f9;",
	}, messages...)
}

type messageComponent struct {
	Message Message
}

func (m *messageComponent) Render(b *web.Builder) error {
	timeStr := m.Message.CreatedAt.Format("15:04:05")

	// Different styles for user and assistant messages
	var bgColor, align, textColor string
	if m.Message.Role == "user" {
		bgColor = "#007bff"
		textColor = "white"
		align = "flex-end"
	} else {
		bgColor = "#e9ecef"
		textColor = "#333"
		align = "flex-start"
	}

	return b.P(
		web.Tag("div", web.Attrs{
			"id":    fmt.Sprintf("message-%d", m.Message.ID),
			"style": fmt.Sprintf("display: flex; justify-content: %s; margin: 10px 0;", align),
		},
			web.Tag("div", web.Attrs{
				"style": fmt.Sprintf("max-width: 70%%; padding: 10px; border-radius: 10px; background: %s; color: %s;", bgColor, textColor),
			},
				web.Tag("div", web.Attrs{"style": "font-weight: bold; margin-bottom: 5px;"},
					Text(m.Message.Role), Text(" #"), Text(strconv.Itoa(m.Message.ID)),
					Text(" - "), Text(timeStr),
				),
				web.Tag("div", web.Attrs{"style": "line-height: 1.4;"},
					m.renderContent(),
				),
			),
		),
	)
}

func (m *messageComponent) renderContent() web.Component {
	if m.Message.Role == "assistant" {
		// Render AI messages as markdown
		return web.Markdown(m.Message.Content)
	} else {
		// Render user messages as plain text with line breaks preserved
		return web.Tag("div", web.Attrs{"style": "white-space: pre-wrap;"},
			web.Esc(m.Message.Content),
		)
	}
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

func (c *chatComponent) renderAutoScrollScript(lastMessageID int) web.Component {
	if lastMessageID == 0 {
		return Text("") // No messages, no need to scroll
	}

	return &autoScrollComponent{LastMessageID: lastMessageID}
}

type autoScrollComponent struct {
	LastMessageID int
}

func (a *autoScrollComponent) Render(b *web.Builder) error {
	script := fmt.Sprintf(`
(function() {
	// Wait for page to fully load, then scroll to last message
	if (document.readyState === 'loading') {
		document.addEventListener('DOMContentLoaded', scrollToLast);
	} else {
		scrollToLast();
	}

	function scrollToLast() {
		var lastMessage = document.getElementById('message-%d');
		if (lastMessage) {
			lastMessage.scrollIntoView({ behavior: 'smooth', block: 'end' });
		}
	}
})();`, a.LastMessageID)

	b.AddInlineScript(script)
	return nil
}
