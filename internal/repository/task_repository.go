package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/DinizJ/desafio/internal/model"
)

//Queries SQL, acesso a banco

// DB

type TaskRepository struct {
	db *sql.DB // connect
}

//Adjust Layers
func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{
		db: db,
	}
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
		return fmt.Errorf("erro ao salvar task no banco:%w", err)
	}

	return nil
}

//FindByID

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

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("erro ao buscar id no banco:%w", err)
	}

	return &task, nil
}

//FindAll

func (r *TaskRepository) FindAll(ctx context.Context, status string) ([]model.Task, error) {
	query := `
 		SELECT id, title, description, status, priority, created_at, updated_at, deleted_at
		FROM tasks
	`

	var (
		rows *sql.Rows
		err  error
	)

	if status != "" {
		query += " WHERE status = ? "
		rows, err = r.db.QueryContext(ctx, query, status)
	} else {
		rows, err = r.db.QueryContext(ctx, query)
	}

	if err != nil {
		return nil, fmt.Errorf("erro ao executar query de tasks:%w", err)
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var task model.Task
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler dados de tasks:%w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao percorrer resultados de tasks:%w", err)
	}

	return tasks, nil
}

//Update

func (r *TaskRepository) Update(ctx context.Context, task *model.Task) error {
	query := `
	UPDATE tasks
	SET title = ?, description = ?, status = ?, priority = ?, updated_at = ?
	WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		task.UpdatedAt,
		task.ID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar tasks:%w", err)
	}
	return nil
}

//Delete

func (r *TaskRepository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM tasks WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar task:%w", err)
	}
	return nil
}
