package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DinizJ/desafio/internal/model"
	"github.com/DinizJ/desafio/internal/repository"
	"github.com/google/uuid"
)

// Lógica de negócio, validações, chamadas ao repository

// Validate

type TaskService struct {
	repo *repository.TaskRepository
}

// Adjust Layers
func NewTaskService(repo *repository.TaskRepository) *TaskService {
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

func TestCreateTask(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		description string
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "valid task",
			title:       "Buy milk",
			description: "2% milk",
			wantErr:     false,
		},
		{
			name:        "empty title",
			title:       "",
			description: "something",
			wantErr:     true,
			errMsg:      "title cannot be empty",
		},
		{
			name:        "title too long",
			title:       string(make([]byte, 256)),
			description: "something",
			wantErr:     true,
			errMsg:      "title too long",
		},
		{
			name:        "valid with empty description",
			title:       "Task",
			description: "",
			wantErr:     false,
		},
	}

	// Mock repository (veremos depois como fazer)
	mockRepo := &mockRepository{}
	service := &TaskService{repo: mockRepo}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := service.CreateTask(context.Background(), tt.title, tt.description)

			//verifica erro
			if tt.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			//Caso sucesso verifica resultado
			if !tt.wantErr {
				if task.Title != tt.title {
					t.Errorf("expected %q, got %q", tt.title, task.Title)
				}
				if task.Status != model.StatusPending {
					t.Errorf("expected status pending, got %q", task.Status)
				}
			}
		})
	}
}

// Teste simples
func TestCompleteTask(t *testing.T) {
	mockRepo := &mockRepository{}
	service := &TaskService{repo: mockRepo}

	//Arrange
	mockRepo.SetupTask(&model.Task{ID: "1", Title: "Test", Status: model.StatusPending})

	//Act
	task, err := service.CompleteTask(context.Background(), "1")

	//Assert
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if task.Status != model.StatusCompleted {
		t.Errorf("expected status completed, got %q", task.Status)
	}
}
