package aiserial

import (
	"encoding/gob"
	"fmt"
	"os"
)

/**
A Basic Serializer that uses a single file and gob encoding
**/

type _serializer struct {
	path string
}

func (s *_serializer) Write(data []SerialData) error {
	file, err := os.Create(s.path)
	if err != nil {
		return fmt.Errorf("Write: %w", err)
	}
	if err := gob.NewEncoder(file).Encode(&data); err != nil {
		return fmt.Errorf("Write: %w", err)
	}
	return nil
}

func (s *_serializer) Read() ([]SerialData, error) {
	file, err := os.Open(s.path)
	if err != nil {
		return nil, fmt.Errorf("Read: %w", err)
	}
	var data []SerialData
	if err := gob.NewDecoder(file).Decode(&data); err != nil {
		return nil, fmt.Errorf("Read: %w", err)
	}
	return data, nil
}

func New(path string) Serializer {
	return &_serializer{path: path}

}
