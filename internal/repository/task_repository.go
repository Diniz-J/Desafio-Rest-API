package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/DinizJ/desafio/internal/model"
)

//Queries SQL, acesso a banco

// DB

type TaskRepository struct {
	db *sql.DB // connect
}

// Save put new task
func (r *TaskRepository) Save(ctx context.Context, task *model.Task) error {
	query := `
		INSERT INTO tasks (id, title, description, status, priority, created_at, updated_at, deleted_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)	    
    `

	_, err := r.db.ExecContext(ctx, query,
		task.ID,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		task.CreatedAt,
		task.UpdatedAt,
		task.DeletedAt,
	)

	if err != nil {
		return fmt.Errorf("falha ao salvar task no banco:%w", err)
	}

	return nil
}

func (r *TaskRepository) FindByID(ctx context.Context, id string) (*model.Task, error) {
	query := `
	SELECT id, title, description, status, priority, created_at, updated_at, deleted_at
 	FROM tasks
 	WHERE id = ?`

	row := r.db.QueryRowContext(ctx, query, id)

	var task model.Task
	err := row.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("falha ao salvar task no banco:%w", err)
	}

	return &task, nil
}
