package ollama

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

const chatEP = "api/chat"

type Generator struct{}

func New() ai.Generator {
	return &Generator{}
}

func (g *Generator) Generate(model, _, baseURL string, conversation ai.Conversation, _ ...ai.Meta) (*ai.GeneratorResponse, error) {
	chatReq := ChatRequest{
		Model: model,
	}
	chatReq.convConv(conversation)

	byts, err := json.MarshalIndent(&chatReq, "", "\t")
	if err != nil {
		return nil, fmt.Errorf("ollama.Generator: %w", err)
	}

	response := ai.GeneratorResponse{}

	start := time.Now()
	resp, err := completion(baseURL, bytes.NewReader(byts))
	if err != nil {
		return nil, fmt.Errorf("ollama.Generator: %w", err)
	}

	response.Elapsed = time.Since(start)
	response.Usage.PromptTokens = int64(resp.PromptEvalCount)
	response.Usage.CompletionTokens = int64(resp.EvalCount)
	response.Usage.TotalTokens = response.Usage.PromptTokens + response.Usage.CompletionTokens
	response.Message.Role = roleAssistant
	response.Message.Text = resp.Message.Content

	return &response, nil
}

func completion(baseURL string, reader io.Reader) (*ChatResponse, error) {
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
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("thread.completion: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("thread.completion: %v %v", resp.StatusCode, resp.Status)
	}

	chatResp := ChatResponse{}
	err = json.NewDecoder(resp.Body).Decode(&chatResp)
	if err != nil {
		return nil, fmt.Errorf("thread.completion: %w", err)
	}
	if len(chatResp.Message.Content) == 0 {
		return nil, fmt.Errorf("thread.completion: No data returned")
	}

	return &chatResp, nil
}
