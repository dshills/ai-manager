package ai

import (
	"fmt"
	"strings"
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

func (m Model) Validate() error {
	var errs []string
	if m.AIName == "" {
		errs = append(errs, "AIName is required")
	}
	if m.Model == "" {
		errs = append(errs, "Model is required")
	}
	if m.BaseURL == "" {
		errs = append(errs, "BaseURL is required")
	}
	if m.Generator == nil {
		errs = append(errs, "Generator is required")
	}
	if len(errs) > 0 {
		return fmt.Errorf("%v", strings.Join(errs, ", "))
	}
	return nil
}

type ThreadData struct {
	ID           string
	AIName       string
	Model        string
	CreatedAt    time.Time
	Conversation aimsg.Conversation
	MetaData     []aimsg.Meta
}
