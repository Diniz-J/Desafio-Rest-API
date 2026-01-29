package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DinizJ/desafio/internal/model"
	"github.com/DinizJ/desafio/internal/repository"
	"github.com/DinizJ/desafio/task"
	"github.com/google/uuid"
)

// Lógica de negócio, validações, chamadas ao repository

// Validate

type TaskService struct {
	repo *repository.TaskRepository
}

func (s *TaskService) CreateTask(ctx context.Context, title string, description string) (*model.Task, error) {

	if title == "" {
		return nil, errors.New("title is required")
	}

	if len(title) > 255 {
		return nil, errors.New("title is too long(max 255)")
	}

	//Cria a TASK
	task := &model.Task{
		ID:          uuid.New().String(), // Gera UUID
		Title:       title,
		Description: description,
		Status:      model.StatusPending,
		Priority:    model.PriorityMedium, // Default
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	//Salva no banco pelo Repository
	if err := s.repo.Save(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) GetTask(ctx context.Context, id string) (*model.Task, error) {

	//Valida se id é uuid valido
	task, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Repo retorna nil se não encontrou
	if task == nil {
		return nil, errors.New("task not found")
	}

	return task, nil
}

// Marca como concluída

func (s *TaskService) CompleteTask(ctx context.Context, id string) (*model.Task, error) {
	task, err := s.GetTask(ctx, id)
	if err != nil {
		return nil, err
	}

	//Valida transição de estado
	if task.Status != model.StatusCompleted {
		return nil, errors.New("task already completed")
	}

	task.Status = model.StatusCompleted
	task.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}
