package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/dshills/ai-manager/aimsg"
)

const chatEP = "api/chat"

func Generator(model, _, baseURL string, conversation aimsg.Conversation, _ ...aimsg.Meta) (aimsg.Message, error) {
	chatReq := ChatRequest{
		Model: model,
	}
	chatReq.convConv(conversation)

	byts, err := json.MarshalIndent(&chatReq, "", "\t")
	if err != nil {
		return aimsg.Message{}, fmt.Errorf("ollama.Generator: %w", err)
	}

	// Make the actual API call
	resp, err := completion(baseURL, bytes.NewReader(byts))
	if err != nil {
		return aimsg.Message{}, fmt.Errorf("ollama.Generator: %w", err)
	}

	msg := aimsg.Message{
		Role: roleAssistant,
		Text: resp.Message.Content,
	}

	return msg, nil
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
