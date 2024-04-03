package gemini

import (
	"ai-manager/aimsg"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const AIName = "gemini"

const (
	baseURL  = "https://generativelanguage.googleapis.com/v1beta"
	geminiEP = "/models/%%MODEL%%:generateContent?key=%%APIKEY%%"
)

func Generator(model, apiKey string, conversation aimsg.Conversation, _ ...aimsg.Meta) (aimsg.Message, error) {
	conlist := []Content{}
	for _, m := range conversation {
		con := Content{Role: m.Role, Parts: []Part{{Text: m.Text}}}
		conlist = append(conlist, con)
	}
	req := Request{Contents: conlist}
	body, err := json.Marshal(&req)
	if err != nil {
		return aimsg.Message{}, fmt.Errorf("gemini.Generator: %w", err)
	}

	resp, err := completion(model, apiKey, bytes.NewReader(body))
	if err != nil {
		return aimsg.Message{}, fmt.Errorf("gemini.Generator: %w", err)
	}

	msg := aimsg.Message{
		Role: resp.Candidates[0].Content.Role,
		Text: resp.Candidates[0].Content.Parts[0].Text,
	}
	return msg, nil
}

func completion(model, apiKey string, reader io.Reader) (*Response, error) {
	ep := fmt.Sprintf("%v%v", baseURL, geminiEP)
	ep = strings.Replace(ep, "%%MODEL%%", model, 1)
	ep = strings.Replace(ep, "%%APIKEY%%", apiKey, 1)

	client := http.Client{}
	req, err := http.NewRequest(http.MethodPost, ep, reader)
	if err != nil {
		return nil, fmt.Errorf("completion: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gemini.completion: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gemini.completion: %v %v", resp.StatusCode, resp.Status)
	}

	chatResp := Response{}
	err = json.NewDecoder(resp.Body).Decode(&chatResp)
	if err != nil {
		return nil, fmt.Errorf("gemini.completion: %w", err)
	}
	if len(chatResp.Candidates) == 0 {
		return nil, fmt.Errorf("gemini.completion: No data returned")
	}
	if len(chatResp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("gemini.completion: No data returned")
	}

	return &chatResp, nil
}
