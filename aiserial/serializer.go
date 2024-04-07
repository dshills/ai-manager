package aiserial

import (
	"encoding/gob"
	"os"

	"github.com/dshills/ai-manager/ai"
)

/**
Serializer is a set of basic functionality to serialize LLM conversations
**/

type Serializer interface {
	Hydrate() ([]ai.ThreadData, error)
	Serialize([]ai.ThreadData) error
}

type _serializer struct {
	path string
}

func (s *_serializer) Serialize(data []ai.ThreadData) error {
	file, err := os.Create(s.path)
	if err != nil {
		return err
	}
	if err := gob.NewEncoder(file).Encode(&data); err != nil {
		return err
	}
	return nil
}

func (s *_serializer) Hydrate() ([]ai.ThreadData, error) {
	file, err := os.Open(s.path)
	if err != nil {
		return nil, err
	}
	var data []ai.ThreadData
	if err := gob.NewDecoder(file).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func New(path string) Serializer {
	return &_serializer{path: path}
}
