package ollama

import (
	"time"

	"github.com/dshills/ai-manager/aigen"
)

const AIName = "ollama"

const (
	roleAssistant = "assistant"
	roleUser      = "user"
)

type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []MessageFrag `json:"messages"`
	Stream   bool          `json:"stream"`
}

func (cr *ChatRequest) convConv(conversation aigen.Conversation) {
	for _, c := range conversation {
		cr.Messages = append(cr.Messages, MessageFrag{Role: c.Role, Content: c.Text})
	}
}

type MessageFrag struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponse struct {
	Model     string    `json:"model"`
	CreatedAt time.Time `json:"created_at"`
	Message   struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	Done               bool  `json:"done"`
	TotalDuration      int64 `json:"total_duration"`
	LoadDuration       int   `json:"load_duration"`
	PromptEvalCount    int   `json:"prompt_eval_count"`
	PromptEvalDuration int   `json:"prompt_eval_duration"`
	EvalCount          int   `json:"eval_count"`
	EvalDuration       int64 `json:"eval_duration"`
}
