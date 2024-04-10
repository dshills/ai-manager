package anthropic

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

const ep = "/messages"

const (
	roleAssistant = "assistant"
	roleUser      = "user"
)

type Generator struct{}

func New() ai.Generator {
	return &Generator{}
}

func (g *Generator) Generate(model, apiKey, baseURL string, conversation ai.Conversation, _ ...ai.Meta) (*ai.GeneratorResponse, error) {
	aireq := Request{Model: model}
	aireq.fillMsgs(conversation)

	body, err := json.Marshal(&aireq)
	if err != nil {
		return nil, fmt.Errorf("anthropic.Generator: %w", err)
	}

	response := ai.GeneratorResponse{}

	start := time.Now()
	resp, err := completion(apiKey, baseURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("anthropic.Generator: %w", err)
	}
	response.Elapsed = time.Since(start)
	response.Usage.PromptTokens = int64(resp.Usage.InputTokens)
	response.Usage.CompletionTokens = int64(resp.Usage.OutputTokens)
	response.Usage.TotalTokens = response.Usage.PromptTokens + response.Usage.CompletionTokens
	response.Message.Role = roleAssistant
	response.Message.Text = resp.Content[0].Text

	return &response, nil
}

func completion(apiKey, baseURL string, reader io.Reader) (*Response, error) {
	ep, err := url.JoinPath(baseURL, ep)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	req, err := http.NewRequest(http.MethodPost, ep, reader)
	if err != nil {
		return nil, fmt.Errorf("completion: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-api-key", apiKey)
	req.Header.Add("anthropic-version", "2023-06-01")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("anthropic.completion: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("anthropic.completion: %v %v", resp.StatusCode, resp.Status)
	}

	chatResp := Response{}
	err = json.NewDecoder(resp.Body).Decode(&chatResp)
	if err != nil {
		return nil, fmt.Errorf("anthropic.completion: %w", err)
	}
	if len(chatResp.Content) == 0 {
		return nil, fmt.Errorf("anthropic.completion: No data returned")
	}
	if len(chatResp.Content[0].Text) == 0 {
		return nil, fmt.Errorf("anthropic.completion: No data returned")
	}

	return &chatResp, nil
}
