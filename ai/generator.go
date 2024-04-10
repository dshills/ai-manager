package ai

import "time"

// GeneratorResponse is the response data from a Generator call
type GeneratorResponse struct {
	Elapsed time.Duration
	Message Message
	Usage   Usage
	Meta    []Meta
}

// Generator is an interface for interacting with an AI
type Generator interface {
	Generate(model, apikey, baseURL string, conversation Conversation, meta ...Meta) (*GeneratorResponse, error)
}
