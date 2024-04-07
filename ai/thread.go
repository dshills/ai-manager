package ai

import (
	"github.com/dshills/ai-manager/aigen"
	"github.com/dshills/ai-manager/aimsg"
)

const (
	ResponseStart    = "--- RESPONSE START ---"
	ResponseComplete = "--- RESPONSE COMPLETE ---"
	ErrorStart       = "--- ERROR START ---"
	ErrorComplete    = "--- ERROR COMPLETE ---"
)

type Thread interface {
	ID() string
	Conversation() aimsg.Conversation
	Info() ThreadData
	Generate(out chan<- string, query string)
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

func (t *_thread) Conversation() aimsg.Conversation {
	return t.info.Conversation
}

func (t *_thread) Info() ThreadData {
	return t.info
}

func (t *_thread) updateConv(msg aimsg.Message) {
	t.info.Conversation = append(t.info.Conversation, msg)
}

func (t *_thread) Generate(out chan<- string, query string) {
	msg := aimsg.Message{Role: "user", Text: query}
	t.updateConv(msg)

	resp, err := t.generator(t.info.Model, t.apiKey, t.baseURL, t.info.Conversation, t.info.MetaData...)
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

func newThread(mgr *Manager, info ThreadData, mod *Model) Thread {
	return &_thread{mgr: mgr, info: info, generator: mod.Generator, apiKey: mod.APIKey, baseURL: mod.BaseURL}
}
