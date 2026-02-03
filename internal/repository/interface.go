package repository

import (
	"context"

	"github.com/DinizJ/desafio/internal/model"
)

// MELHORIA: Interface criada para permitir testes unitários com mocks
// Agora o service pode receber qualquer implementação que satisfaça essa interface,
// seja o repository real (com banco) ou um mock (para testes)

type TaskRepositoryInterface interface {
	Save(ctx context.Context, task *model.Task) error
	FindByID(ctx context.Context, id string) (*model.Task, error)
	FindAll(ctx context.Context, status string) ([]model.Task, error)
	Update(ctx context.Context, task *model.Task) error
	Delete(ctx context.Context, id string) error
}

// Verifica em tempo de compilação se TaskRepository implementa a interface
var _ TaskRepositoryInterface = (*TaskRepository)(nil)
