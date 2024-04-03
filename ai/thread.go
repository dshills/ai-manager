package ai

import (
	"ai-manager/aigen"
	"ai-manager/aimsg"
	"ai-manager/aiserial"
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
	Info() aiserial.SerialData
	Generate(out chan<- string, query string)
}

type _thread struct {
	info      aiserial.SerialData
	generator aigen.Generator
	apiKey    string
	mgr       *Manager
}

func (t *_thread) ID() string {
	return t.info.ID
}

func (t *_thread) Conversation() aimsg.Conversation {
	return t.info.Conversation
}

func (t *_thread) Info() aiserial.SerialData {
	return t.info
}

func (t *_thread) updateConv(msg aimsg.Message) {
	t.info.Conversation = append(t.info.Conversation, msg)
}

func (t *_thread) Generate(out chan<- string, query string) {
	msg := aimsg.Message{Role: "user", Text: query}
	t.updateConv(msg)
	resp, err := t.generator(t.info.Model, t.apiKey, t.info.Conversation, t.info.MetaData...)
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

func newThread(mgr *Manager, info aiserial.SerialData, gen aiGenerator) Thread {
	return &_thread{mgr: mgr, info: info, generator: gen.Generator, apiKey: gen.APIKey}
}
