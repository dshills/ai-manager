package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dshills/ai-manager/aimsg"
)

const chatEP = "https://api.openai.com/v1/chat/completions"

const AIName = "openai"

const (
	roleSystem    = "system"
	roleAssistant = "assistant"
	roleUser      = "user"
)

func Generator(model, apiKey string, conversation aimsg.Conversation, _ ...aimsg.Meta) (aimsg.Message, error) {
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
		return aimsg.Message{}, fmt.Errorf("openai.Generator: %w", err)
	}

	// Make the actual API call
	resp, err := completion(apiKey, bytes.NewReader(byts))
	if err != nil {
		return aimsg.Message{}, fmt.Errorf("openai.Generator: %w", err)
	}

	msg := aimsg.Message{
		Role: roleAssistant,
		Text: resp.Choices[0].Message.Content,
	}

	return msg, nil
}

func completion(apiKey string, reader io.Reader) (*ChatResp, error) {
	client := http.Client{}
	req, err := http.NewRequest(http.MethodPost, chatEP, reader)
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
