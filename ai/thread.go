package ai

import (
	"github.com/dshills/ai-manager/aigen"
)

const (
	ResponseStart    = "--- RESPONSE START ---"
	ResponseComplete = "--- RESPONSE COMPLETE ---"
	ErrorStart       = "--- ERROR START ---"
	ErrorComplete    = "--- ERROR COMPLETE ---"
)

type Thread interface {
	ID() string
	Conversation() aigen.Conversation
	Info() ThreadData
	Generate(out chan<- string, query string)
	Converse(query string) (string, error)
}

type _thread struct {
	info      ThreadData
	generator aigen.Generator
	apiKey    string
	mgr       *Manager
	baseURL   string
}

func (t *_thread) ID() string {
	return t.info.ID
}

func (t *_thread) Conversation() aigen.Conversation {
	return t.info.Conversation
}

func (t *_thread) Info() ThreadData {
	return t.info
}

func (t *_thread) updateConv(msg aigen.Message) {
	t.info.Conversation = append(t.info.Conversation, msg)
}

func (t *_thread) updateUsage(usage aigen.Usage) {
	t.info.PromptTokens += usage.PromptTokens
	t.info.CompletionTokens += usage.CompletionTokens
	t.info.TotalTokens += usage.TotalTokens
}

func (t *_thread) Generate(out chan<- string, query string) {
	msg := aigen.Message{Role: "user", Text: query}
	t.updateConv(msg)

	resp, usage, err := t.generator(t.info.Model, t.apiKey, t.baseURL, t.info.Conversation, t.info.MetaData...)
	t.updateUsage(usage) // Even if error we may have had usage tokens
	if err != nil {
		out <- ErrorStart
		out <- err.Error()
		out <- ErrorComplete
		return
	}
	t.updateConv(resp)
	out <- ResponseStart
	out <- resp.Text
	out <- ResponseComplete
}

func (t *_thread) Converse(query string) (string, error) {
	msg := aigen.Message{Role: "user", Text: query}
	t.updateConv(msg)

	resp, usage, err := t.generator(t.info.Model, t.apiKey, t.baseURL, t.info.Conversation, t.info.MetaData...)
	t.updateUsage(usage) // Even if error we may have had usage tokens
	if err != nil {
		return "", err
	}
	t.updateConv(resp)
	return resp.Text, nil
}

func newThread(mgr *Manager, info ThreadData, mod *Model) Thread {
	return &_thread{mgr: mgr, info: info, generator: mod.Generator, apiKey: mod.APIKey, baseURL: mod.BaseURL}
}
