package task

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type Store interface {
	Load() ([]Task, error)
	Save([]Task) error
}

type JSONStore struct {
	path string
	mu   sync.Mutex
}

func NewJSONStore(path string) *JSONStore {
	return &JSONStore{path: path}
}

func (s *JSONStore) Load() ([]Task, error) {
	s.mu.Lock()

	defer s.mu.Unlock()

	if _, err := os.Stat(s.path); errors.Is(err, os.ErrNotExist) {
		return []Task{}, nil
	}

	data, err := os.ReadFile(s.path)

	if err != nil {
		return nil, err
	}

	var tasks []Task

	if len(data) == 0 {
		return []Task{}, nil
	}

	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *JSONStore) Save(tasks []Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.MarshalIndent(tasks, "", " ")

	if err != nil {
		return err
	}

	if err := os.MkdirAll("data", 0o755); err != nil {
		return err
	}

	return os.WriteFile(s.path, data, 0o644)
}
