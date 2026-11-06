package task

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrTaskNotFound     = errors.New("task not found")
	ErrEmptyDescription = errors.New("description cannot be empty")
	ErrInvalidStatus    = errors.New("invalid status")
)

type Storage interface {
	Load() ([]Task, error)
	Save(tasks []Task) error
}

type Service struct {
	storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) Add(description string) (*Task, error) {
	trimmed := strings.TrimSpace(description)
	if trimmed == "" {
		return nil, ErrEmptyDescription
	}

	tasks, err := s.storage.Load()
	if err != nil {
		return nil, err
	}

	nextID := 1
	for _, t := range tasks {
		if t.ID >= nextID {
			nextID = t.ID + 1
		}
	}

	now := time.Now()
	newTask := Task{
		ID:          nextID,
		Description: trimmed,
		Status:      StatusTodo,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tasks = append(tasks, newTask)
	if err := s.storage.Save(tasks); err != nil {
		return nil, err
	}

	return &newTask, nil
}

func (s *Service) Update(id int, description string) error {
	trimmed := strings.TrimSpace(description)
	if trimmed == "" {
		return ErrEmptyDescription
	}

	tasks, err := s.storage.Load()
	if err != nil {
		return err
	}

	idx := findTaskIndex(tasks, id)
	if idx == -1 {
		return ErrTaskNotFound
	}

	tasks[idx].Description = trimmed
	tasks[idx].UpdatedAt = time.Now()

	return s.storage.Save(tasks)
}

func (s *Service) Delete(id int) error {
	tasks, err := s.storage.Load()
	if err != nil {
		return err
	}

	idx := findTaskIndex(tasks, id)
	if idx == -1 {
		return ErrTaskNotFound
	}

	tasks = append(tasks[:idx], tasks[idx+1:]...)

	return s.storage.Save(tasks)
}

func (s *Service) MarkStatus(id int, status Status) error {
	if !status.IsValid() {
		return ErrInvalidStatus
	}

	tasks, err := s.storage.Load()
	if err != nil {
		return err
	}

	idx := findTaskIndex(tasks, id)
	if idx == -1 {
		return ErrTaskNotFound
	}

	tasks[idx].Status = status
	tasks[idx].UpdatedAt = time.Now()

	return s.storage.Save(tasks)
}

func (s *Service) List(statusFilter Status) ([]Task, error) {
	tasks, err := s.storage.Load()
	if err != nil {
		return nil, err
	}

	if statusFilter == "" {
		return tasks, nil
	}

	if !statusFilter.IsValid() {
		return nil, ErrInvalidStatus
	}

	filtered := make([]Task, 0)
	for _, t := range tasks {
		if t.Status == statusFilter {
			filtered = append(filtered, t)
		}
	}

	return filtered, nil
}

func findTaskIndex(tasks []Task, id int) int {
	for i, t := range tasks {
		if t.ID == id {
			return i
		}
	}
	return -1
}
