package service

import (
	"context"
	"testing"
	"time"

	"github.com/DinizJ/desafio/internal/model"
)

// CORREÇÃO: Testes movidos para arquivo separado (_test.go)
// Isso é a convenção padrão do Go para testes unitários.

// Mock Repository - Implementação fake para testes
// Simula o comportamento do repository sem acessar banco de dados real
type mockRepository struct {
	tasks map[string]*model.Task
}

// Save simula salvar uma task no banco
func (m *mockRepository) Save(ctx context.Context, task *model.Task) error {
	if m.tasks == nil {
		m.tasks = make(map[string]*model.Task)
	}
	m.tasks[task.ID] = task
	return nil
}

// FindByID simula buscar uma task por ID
func (m *mockRepository) FindByID(ctx context.Context, id string) (*model.Task, error) {
	task, ok := m.tasks[id]
	if !ok {
		return nil, nil
	}
	return task, nil
}

// FindAll simula listar todas as tasks
func (m *mockRepository) FindAll(ctx context.Context, status string) ([]model.Task, error) {
	var result []model.Task
	for _, task := range m.tasks {
		if status == "" || task.Status == status {
			result = append(result, *task)
		}
	}
	return result, nil
}

// Update simula atualizar uma task
func (m *mockRepository) Update(ctx context.Context, task *model.Task) error {
	if m.tasks == nil {
		m.tasks = make(map[string]*model.Task)
	}
	m.tasks[task.ID] = task
	return nil
}

// Delete simula deletar uma task
func (m *mockRepository) Delete(ctx context.Context, id string) error {
	delete(m.tasks, id)
	return nil
}

// SetupTask é um helper para configurar o mock com dados de teste
func (m *mockRepository) SetupTask(task *model.Task) {
	if m.tasks == nil {
		m.tasks = make(map[string]*model.Task)
	}
	m.tasks[task.ID] = task
}

// ------------------------ TESTES ------------------------

func TestCreateTask(t *testing.T) {
	// Table-driven tests: cada caso de teste em uma estrutura
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
			errMsg:      "title is required",
		},
		{
			name:        "title too long",
			title:       string(make([]byte, 256)),
			description: "something",
			wantErr:     true,
			errMsg:      "title is too long",
		},
		{
			name:        "valid with empty description",
			title:       "Task",
			description: "",
			wantErr:     false,
		},
	}

	// Cria mock repository e service
	mockRepo := &mockRepository{}
	service := &TaskService{repo: mockRepo}

	// Executa cada caso de teste
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := service.CreateTask(context.Background(), tt.title, tt.description)

			// Verifica se erro ocorreu quando esperado
			if tt.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Caso sucesso, verifica resultado
			if !tt.wantErr {
				if task.Title != tt.title {
					t.Errorf("expected title %q, got %q", tt.title, task.Title)
				}
				if task.Status != model.StatusPending {
					t.Errorf("expected status pending, got %q", task.Status)
				}
				if task.Priority != model.PriorityMedium {
					t.Errorf("expected default priority medium, got %q", task.Priority)
				}
			}
		})
	}
}

func TestCompleteTask(t *testing.T) {
	mockRepo := &mockRepository{}
	service := &TaskService{repo: mockRepo}

	// Arrange: Configura uma task pending
	mockRepo.SetupTask(&model.Task{
		ID:        "1",
		Title:     "Test Task",
		Status:    model.StatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	// Act: Marca como concluída
	task, err := service.CompleteTask(context.Background(), "1")

	// Assert: Verifica resultado
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if task.Status != model.StatusCompleted {
		t.Errorf("expected status completed, got %q", task.Status)
	}
}

func TestCompleteTask_AlreadyCompleted(t *testing.T) {
	mockRepo := &mockRepository{}
	service := &TaskService{repo: mockRepo}

	// Arrange: Task já está completed
	mockRepo.SetupTask(&model.Task{
		ID:        "1",
		Title:     "Test Task",
		Status:    model.StatusCompleted,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	// Act: Tenta marcar como concluída novamente
	task, err := service.CompleteTask(context.Background(), "1")

	// Assert: Deve retornar erro
	if err == nil {
		t.Error("expected error for already completed task, got nil")
	}
	if task != nil {
		t.Errorf("expected nil task, got %v", task)
	}
}

func TestUpdateTask_ValidateEnums(t *testing.T) {
	mockRepo := &mockRepository{}
	service := &TaskService{repo: mockRepo}

	// Setup: Cria uma task
	mockRepo.SetupTask(&model.Task{
		ID:        "1",
		Title:     "Original",
		Status:    model.StatusPending,
		Priority:  model.PriorityMedium,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	tests := []struct {
		name     string
		status   string
		priority string
		wantErr  bool
	}{
		{
			name:     "valid status and priority",
			status:   model.StatusCompleted,
			priority: model.PriorityHigh,
			wantErr:  false,
		},
		{
			name:     "invalid status",
			status:   "invalid_status",
			priority: model.PriorityLow,
			wantErr:  true,
		},
		{
			name:     "invalid priority",
			status:   model.StatusPending,
			priority: "invalid_priority",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.UpdateTask(context.Background(), "1", "Title", "Desc", tt.status, tt.priority)

			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
