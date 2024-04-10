package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/dshills/ai-manager/ai"
)

const chatEP = "/chat/completions"

const AIName = "openai"

const (
	roleSystem    = "system"
	roleAssistant = "assistant"
	roleUser      = "user"
)

type Generator struct{}

func New() ai.Generator {
	return &Generator{}
}

func (g *Generator) Generate(model, apiKey, baseURL string, conversation ai.Conversation, _ ...ai.Meta) (*ai.GeneratorResponse, error) {
	frags := []MessageFrag{}
	for _, m := range conversation {
		frags = append(frags, MessageFrag{Role: m.Role, Content: m.Text})
	}
	chatReq := CreateRequest{
		Model:    model,
		Messages: frags,
	}
	byts, err := json.MarshalIndent(&chatReq, "", "\t")
	if err != nil {
		return nil, fmt.Errorf("openai.Generator: %w", err)
	}

	response := ai.GeneratorResponse{}

	start := time.Now()
	resp, err := completion(apiKey, baseURL, bytes.NewReader(byts))
	if err != nil {
		return nil, fmt.Errorf("openai.Generator: %w", err)
	}

	response.Elapsed = time.Since(start)
	response.Usage.PromptTokens = resp.Usage.PromptTokens
	response.Usage.CompletionTokens = resp.Usage.CompletionTokens
	response.Usage.TotalTokens = resp.Usage.TotalTokens
	response.Message.Role = roleAssistant
	response.Message.Text = resp.Choices[0].Message.Content

	return &response, nil
}

func completion(apiKey, baseURL string, reader io.Reader) (*ChatResp, error) {
	ep, err := url.JoinPath(baseURL, chatEP)
	if err != nil {
		return nil, err
	}
	client := http.Client{}
	req, err := http.NewRequest(http.MethodPost, ep, reader)
	if err != nil {
		return nil, fmt.Errorf("thread.completion: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("thread.completion: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("thread.completion: %v %v", resp.StatusCode, resp.Status)
	}

	chatResp := ChatResp{}
	err = json.NewDecoder(resp.Body).Decode(&chatResp)
	if err != nil {
		return nil, fmt.Errorf("thread.completion: %w", err)
	}
	if len(chatResp.Choices) == 0 {
		return nil, fmt.Errorf("thread.completion: No data returned")
	}

	return &chatResp, nil
}
