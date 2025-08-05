package infrastructure

import (
	"context"
	"errors"
	"fmt"

	"github.com/Mayer-04/grpc-task-manager-go/internal/tasks/domain"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepositoryImpl struct {
	dbpool *pgxpool.Pool
}

func NewTaskRepository(dbpool *pgxpool.Pool) domain.TaskRepository {
	return &TaskRepositoryImpl{
		dbpool: dbpool,
	}
}

func (t *TaskRepositoryImpl) CreateTask(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	// Generate a new UUID
	taskID, err := uuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("failed to generate UUID: %w", err)
	}

	// SQL query to insert a new task
	query := `
		INSERT INTO tasks (id, user_id, title, description, completed) 
			VALUES ($1, $2, $3, $4, $5) 
		RETURNING id, user_id, title, description, completed, created_at, updated_at;`

	// Task instance to hold the result
	newTask := &domain.Task{}

	// Execute the query and scan the result into newTask
	err = t.dbpool.QueryRow(
		ctx,
		query,
		taskID,
		task.UserID,
		task.Title,
		task.Description,
		task.Completed,
	).Scan(
		&newTask.ID,
		&newTask.UserID,
		&newTask.Title,
		&newTask.Description,
		&newTask.Completed,
		&newTask.CreatedAt,
		&newTask.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to insert task: %w", err)
	}

	return newTask, nil
}

func (t *TaskRepositoryImpl) GetTask(ctx context.Context, id string) (*domain.Task, error) {
	query := `
		SELECT id, user_id, title, description, completed, created_at, updated_at 
		FROM tasks 
		WHERE id = $1;`

	task := &domain.Task{}

	err := t.dbpool.QueryRow(ctx, query, id).Scan(
		&task.ID,
		&task.UserID,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("task not found: %w", err)
		}
		return nil, fmt.Errorf("failed to retrieve task: %w", err)
	}

	return task, nil
}

func (t *TaskRepositoryImpl) UpdateTask(ctx context.Context, task *domain.Task) (*domain.Task, error) {

	query := `
		UPDATE tasks 
		SET title = $1, description = $2, completed = $3, updated_at = NOW() 
		WHERE id = $4 
		RETURNING id, user_id, title, description, completed, created_at, updated_at;`

	newTask := &domain.Task{}

	err := t.dbpool.QueryRow(ctx, query, task.Title, task.Description, task.Completed, task.ID).Scan(
		&newTask.ID,
		&newTask.UserID,
		&newTask.Title,
		&newTask.Description,
		&newTask.Completed,
		&newTask.CreatedAt,
		&newTask.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	return newTask, nil
}

func (r *TaskRepositoryImpl) DeleteTask(ctx context.Context, id string) error {
	const query = "DELETE FROM tasks WHERE id = $1;"

	result, err := r.dbpool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("could not delete task: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func (t *TaskRepositoryImpl) ListTasksByUser(ctx context.Context, userID string) ([]*domain.Task, error) {

	query := `
		SELECT id, user_id, title, description, completed, created_at, updated_at 
		FROM tasks 
		WHERE user_id = $1;`

	rows, err := t.dbpool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	defer rows.Close()

	var tasks []*domain.Task
	for rows.Next() {
		task := &domain.Task{}
		err := rows.Scan(
			&task.ID,
			&task.UserID,
			&task.Title,
			&task.Description,
			&task.Completed,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tasks: %w", err)
	}

	return tasks, nil
}

func (t *TaskRepositoryImpl) MarkTaskComplete(ctx context.Context, taskID string) (*domain.Task, error) {
	const query = `
		UPDATE tasks 
		SET completed = true, updated_at = NOW() 
		WHERE id = $1 
		RETURNING id, user_id, title, description, completed, created_at, updated_at;
	`

	var updatedTask domain.Task
	err := t.dbpool.QueryRow(ctx, query, taskID).Scan(
		&updatedTask.ID,
		&updatedTask.UserID,
		&updatedTask.Title,
		&updatedTask.Description,
		&updatedTask.Completed,
		&updatedTask.CreatedAt,
		&updatedTask.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to mark task complete: %w", err)
	}

	return &updatedTask, nil
}
