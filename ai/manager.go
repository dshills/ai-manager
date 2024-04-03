package ai

import (
	"ai-manager/aigen"
	"ai-manager/aimsg"
	"ai-manager/aiserial"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type aiGenerator struct {
	AIName    string
	APIKey    string
	Generator aigen.Generator
	Models    []string
}

type Manager struct {
	threads       []Thread
	currentThread Thread
	generators    map[string]aiGenerator
	serializer    aiserial.Serializer
	serSync       sync.Mutex
}

func (ai *Manager) Models() []string {
	models := []string{}
	for _, aig := range ai.generators {
		for _, model := range aig.Models {
			models = append(models, fmt.Sprintf("%v: %v", aig.AIName, model))
		}
	}
	sort.Strings(models)
	return models
}

func (ai *Manager) RemoveThread(threadID string) error {
	idx := -1
	for i, thread := range ai.threads {
		if thread.ID() == threadID {
			idx = i
			break
		}
	}
	if idx == -1 {
		return fmt.Errorf("RemoveThread: Not found")
	}
	ai.threads = append(ai.threads[:idx], ai.threads[idx+1:]...)
	return ai.Serialize()
}

func (ai *Manager) Hydrate() error {
	ai.serSync.Lock()
	defer ai.serSync.Unlock()

	if ai.serializer == nil {
		return fmt.Errorf("Hydrate: no defined serializer")
	}

	sdata, err := ai.serializer.Read()
	if err != nil {
		return fmt.Errorf("Hydrate: %w", err)
	}
	var errs []string

	for i := range sdata {
		gen, err := ai.generatorInfo(sdata[i].AIName)
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		thread := newThread(ai, sdata[i], gen)
		ai.threads = append(ai.threads, thread)
	}
	if len(ai.threads) > 0 && ai.currentThread == nil {
		ai.currentThread = ai.threads[len(ai.threads)-1]
	}
	if len(errs) > 0 {
		return fmt.Errorf("Hydrate: %v", strings.Join(errs, ", "))
	}

	return nil
}

func (ai *Manager) Serialize() error {
	ai.serSync.Lock()
	defer ai.serSync.Unlock()

	if ai.serializer == nil {
		return fmt.Errorf("Serialize: no defined serializer")
	}
	convs := make([]aiserial.SerialData, len(ai.threads))
	for _, thread := range ai.threads {
		convs = append(convs, thread.Info())
	}
	if err := ai.serializer.Write(convs); err != nil {
		return fmt.Errorf("Serialize: %w", err)
	}

	return nil
}

func (ai *Manager) Conversations() []aiserial.SerialData {
	if len(ai.threads) == 0 {
		return nil
	}
	var convs []aiserial.SerialData
	for _, thread := range ai.threads {
		dis := thread.Info()
		convs = append(convs, dis)
	}
	return convs
}

func (ai *Manager) RegisterGenerator(aiName, apiKey string, models []string, generator aigen.Generator) {
	ai.generators[strings.ToLower(aiName)] = aiGenerator{AIName: aiName, APIKey: apiKey, Generator: generator, Models: models}
}

func (ai *Manager) generatorInfo(aiName string) (aiGenerator, error) {
	aig, ok := ai.generators[strings.ToLower(aiName)]
	if !ok {
		return aiGenerator{}, fmt.Errorf("%v Generator not found", aiName)
	}
	return aig, nil
}

func (ai *Manager) NewThread(aiName, model string, meta ...aimsg.Meta) error {
	gen, err := ai.generatorInfo(aiName)
	if err != nil {
		return err
	}
	info := aiserial.SerialData{
		ID:        uuid.New().String(),
		AIName:    aiName,
		Model:     model,
		CreatedAt: time.Now(),
		MetaData:  meta,
	}
	thread := newThread(ai, info, gen)
	ai.threads = append(ai.threads, thread)
	ai.currentThread = thread

	return nil
}

func (ai *Manager) SwitchThread(threadID string) error {
	for i := range ai.threads {
		if ai.threads[i].ID() == threadID {
			ai.currentThread = ai.threads[i]
			return nil
		}
	}
	return fmt.Errorf("SwitchThread: thread not found")
}

func (ai *Manager) CurrentThread() Thread {
	return ai.currentThread
}

func New(serializer aiserial.Serializer) *Manager {
	aigen := Manager{
		generators: make(map[string]aiGenerator),
		serializer: serializer,
	}
	return &aigen
}
