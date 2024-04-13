package ai

import (
	"time"

	"github.com/dshills/ai-manager/aitool"
)

// GeneratorResponse is the response data from a Generator call
type GeneratorResponse struct {
	Elapsed      time.Duration
	Message      Message
	Usage        Usage
	Meta         []Meta
	ToolCalls    []ToolCall
	FinishReason string
}

type ToolCall struct {
	ID   string
	Type string
	Name string
	Args string
}

// Generator is an interface for interacting with an AI
type Generator interface {
	Model() string
	Generate(conversation Conversation, meta []Meta, tools []aitool.Tool) (*GeneratorResponse, error)
}
