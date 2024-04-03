package mistral

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dshills/ai-manager/aimsg"
)

const mistralEP = "https://api.mistral.ai/v1/chat/completions"

const AIName = "mistral"

func Generator(model, apiKey string, conversation aimsg.Conversation, _ ...aimsg.Meta) (aimsg.Message, error) {
	messages := []Message{}
	for _, m := range conversation {
		msg := Message{Role: m.Role, Content: m.Text}
		messages = append(messages, msg)
	}
	req := Request{
		Model:       model,
		Messages:    messages,
		Stream:      false,
		SafePrompt:  false,
		Temperature: 0.2,
	}
	body, err := json.Marshal(&req)
	if err != nil {
		return aimsg.Message{}, fmt.Errorf("mistral.Generator: %w", err)
	}

	resp, err := completion(apiKey, bytes.NewReader(body))
	if err != nil {
		return aimsg.Message{}, fmt.Errorf("mistral.Generator: %w", err)
	}

	msg := aimsg.Message{
		Role: resp.Choices[0].Message.Role,
		Text: resp.Choices[0].Message.Content,
	}
	return msg, nil
}

func completion(apiKey string, reader io.Reader) (*Response, error) {
	client := http.Client{}
	req, err := http.NewRequest(http.MethodPost, mistralEP, reader)
	if err != nil {
		return nil, fmt.Errorf("completion: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("mistral.completion: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("mistral.completion: %v %v", resp.StatusCode, resp.Status)
	}

	chatResp := Response{}
	err = json.NewDecoder(resp.Body).Decode(&chatResp)
	if err != nil {
		return nil, fmt.Errorf("mistral.completion: %w", err)
	}
	if len(chatResp.Choices) == 0 {
		return nil, fmt.Errorf("mistral.completion: No data returned")
	}
	if len(chatResp.Choices[0].Message.Content) == 0 {
		return nil, fmt.Errorf("mistral.completion: No data returned")
	}

	return &chatResp, nil
}
