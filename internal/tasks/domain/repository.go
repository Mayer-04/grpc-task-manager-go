package domain

import "context"

type TaskRepository interface {
	CreateTask(ctx context.Context, task *Task) (*Task, error)
	GetTask(ctx context.Context, id string) (*Task, error)
	UpdateTask(ctx context.Context, task *Task) (*Task, error)
	DeleteTask(ctx context.Context, id string) error
	ListTasksByUser(ctx context.Context, userID string) ([]*Task, error)
	MarkTaskComplete(ctx context.Context, id string) (*Task, error)
}
