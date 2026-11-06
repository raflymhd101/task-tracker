package storage

import (
	"encoding/json"
	"errors"
	"os"
	"task-tracker/internal/task"
)

type Storage struct {
	filePath string
}

func New(filePath string) *Storage {
	return &Storage{filePath: filePath}
}

func (s *Storage) Load() ([]task.Task, error) {
	data, err := os.ReadFile(s.filePath)

	if errors.Is(err, os.ErrNotExist) {
		return []task.Task{}, nil
	}

	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return []task.Task{}, nil
	}

	var tasks []task.Task

	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *Storage) Save(task []task.Task) error {
	data, err := json.MarshalIndent(task, "", "  ")

	if err != nil {
		return err
	}

	return os.WriteFile(s.filePath, data, 0644)
}
