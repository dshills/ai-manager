package ai

import (
	"fmt"
	"strings"
	"sync"

	"github.com/google/uuid"
)

type Manager struct {
	threads       []Thread
	currentThread Thread
	models        map[string]Model
	m             sync.RWMutex
	modelNames    []string
}

func New() *Manager {
	aigen := Manager{models: make(map[string]Model)}
	return &aigen
}

func (ai *Manager) NewThread(info ThreadData) (Thread, error) {
	gen, err := ai.generatorInfo(info.AIName, info.Model)
	if err != nil {
		return nil, err
	}
	if info.ID == "" {
		info.ID = uuid.New().String()
	}
	thread := newThread(ai, info, gen)
	ai.threads = append(ai.threads, thread)
	ai.currentThread = thread

	return thread, nil
}

func (ai *Manager) SwitchThread(threadID string) (Thread, error) {
	for i := range ai.threads {
		if ai.threads[i].ID() == threadID {
			ai.currentThread = ai.threads[i]
			return ai.threads[i], nil
		}
	}
	return nil, fmt.Errorf("SwitchThread: thread not found")
}

func (ai *Manager) CurrentThread() Thread {
	return ai.currentThread
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
	return nil
}

func (ai *Manager) Threads() []ThreadData {
	if len(ai.threads) == 0 {
		return []ThreadData{}
	}
	convs := make([]ThreadData, len(ai.threads))
	for i, thread := range ai.threads {
		convs[i] = thread.Info()
	}
	return convs
}

func (ai *Manager) RegisterGenerators(models ...Model) error {
	ai.m.Lock()
	defer ai.m.Unlock()

	for _, mod := range models {
		if err := mod.Validate(); err != nil {
			return err
		}
		key := fmt.Sprintf("%v:%v", strings.ToLower(mod.AIName), strings.ToLower(mod.Model))
		ai.models[key] = mod
	}
	names := []string{}
	for key := range ai.models {
		names = append(names, key)
	}
	ai.modelNames = names
	return nil
}

func (ai *Manager) Models() []string {
	return ai.modelNames
}

func (ai *Manager) generatorInfo(aiName, model string) (*Model, error) {
	ai.m.RLock()
	defer ai.m.RUnlock()

	key := fmt.Sprintf("%v:%v", strings.ToLower(aiName), strings.ToLower(model))
	mod, ok := ai.models[key]
	if !ok {
		return nil, fmt.Errorf("%v %v Generator not found", aiName, model)
	}
	return &mod, nil
}
