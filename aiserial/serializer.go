package aiserial

import (
	"ai-manager/aimsg"
	"time"
)

/**
Serializer is a set of functionality to serialize LLM conversations
**/

type SerialData struct {
	ID           string
	AIName       string
	Model        string
	CreatedAt    time.Time
	Conversation aimsg.Conversation
	MetaData     []aimsg.Meta
}

type Serializer interface {
	Write([]SerialData) error
	Read() ([]SerialData, error)
}
