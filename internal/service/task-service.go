package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/DinizJ/desafio/internal/model"
	"github.com/DinizJ/desafio/internal/repository"
	"github.com/google/uuid"
)

// Lógica de negócio, validações, chamadas ao repository

// MELHORIA: Service agora usa interface em vez de tipo concreto
// Isso facilita testes unitários com mocks e torna o código mais flexível

type TaskService struct {
	repo repository.TaskRepositoryInterface
}

// ------------------------CREATE TASK--------------------------------
// Adjust Layers
func NewTaskService(repo repository.TaskRepositoryInterface) *TaskService {
	return &TaskService{repo: repo}
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

// ------------------------GET TASK--------------------------------
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

// ------------------------COMPLETE TASK--------------------------------
func (s *TaskService) CompleteTask(ctx context.Context, id string) (*model.Task, error) {
	task, err := s.GetTask(ctx, id)
	if err != nil {
		return nil, err
	}

	// CORREÇÃO: A lógica estava invertida. Se o status JÁ É "completed",
	// devemos retornar erro. A condição anterior verificava se era DIFERENTE (!= ),
	// o que causava o comportamento oposto ao esperado.
	if task.Status == model.StatusCompleted {
		return nil, errors.New("task already completed")
	}

	task.Status = model.StatusCompleted
	task.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

// ------------------------DELETE TASK--------------------------------
func (s *TaskService) DeleteTask(ctx context.Context, id string) error {

	task, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if task == nil {
		return errors.New("task not found")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

// ------------------------UPDATE TASK--------------------------------
func (s *TaskService) UpdateTask(
	ctx context.Context, id string, title string, description string, status string, priority string,
) (*model.Task, error) {
	task, err := s.GetTask(ctx, id)
	if err != nil {
		return nil, err
	}

	if task == nil {
		return nil, errors.New("task not found")
	}

	// MELHORIA: Validar title se fornecido
	if title != "" {
		if len(title) > 255 {
			return nil, errors.New("title is too long (max 255)")
		}
		task.Title = title
	}

	if description != "" {
		task.Description = description
	}

	// CORREÇÃO: Validar status enum se fornecido
	// Antes aceitava qualquer valor, agora valida contra as constantes do model
	if status != "" {
		if status != model.StatusPending && status != model.StatusCompleted {
			return nil, errors.New("invalid status: must be 'pending' or 'completed'")
		}
		task.Status = status
	}

	// CORREÇÃO: Validar priority enum se fornecido
	// Antes aceitava qualquer valor, agora valida contra as constantes do model
	if priority != "" {
		if priority != model.PriorityLow && priority != model.PriorityMedium && priority != model.PriorityHigh {
			return nil, errors.New("invalid priority: must be 'low', 'medium' or 'high'")
		}
		task.Priority = priority
	}

	task.UpdatedAt = time.Now()

	err = s.repo.Update(ctx, task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

// ------------------------LIST TASK--------------------------------
func (s *TaskService) ListTask(ctx context.Context, status string) ([]model.Task, error) {

	tasks, err := s.repo.FindAll(ctx, status)
	if err != nil {
		return nil, fmt.Errorf("Error listing tasks: %w", err)
	}
	return tasks, nil
}
