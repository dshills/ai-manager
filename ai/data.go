package ai

import (
	"time"

	"github.com/dshills/ai-manager/aigen"
	"github.com/dshills/ai-manager/aimsg"
)

type Model struct {
	AIName    string
	Model     string
	APIKey    string
	BaseURL   string
	Generator aigen.Generator
}

type ThreadData struct {
	ID           string
	AIName       string
	Model        string
	CreatedAt    time.Time
	Conversation aimsg.Conversation
	MetaData     []aimsg.Meta
}
