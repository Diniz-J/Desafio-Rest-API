package service

import (
	"fmt"

	"github.com/DinizJ/desafio/task"
	"github.com/google/uuid"
)

// Lógica de negócio, validações, chamadas ao repository

// Validate

type TaskService struct {
	//repository
}

func (s *TaskService) CreateId() (string, error) {

	id := uuid.New().String()

	t := &Task{
		ID: id,
	}
	if err := s.repo.Save(t); err != nil {
		return "", fmt.Errorf("Cannot create task: %w", err)
	}

	return id, nil
}

func (s *TaskService) CreateTitle()
