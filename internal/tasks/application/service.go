package application

import (
	"context"
	"fmt"

	"github.com/Mayer-04/grpc-task-manager-go/internal/tasks/domain"
	"github.com/gofrs/uuid"
)

type TaskService struct {
	taskRepo domain.TaskRepository
}

func NewTaskService(taskRepo domain.TaskRepository) *TaskService {
	return &TaskService{
		taskRepo: taskRepo,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, userID, title, description string, completed bool) (*domain.Task, error) {
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}
	if title == "" {
		return nil, fmt.Errorf("title is required")
	}

	task := &domain.Task{
		UserID:      userID,
		Title:       title,
		Description: description,
		Completed:   completed,
	}

	return s.taskRepo.CreateTask(ctx, task)
}

func (s *TaskService) GetTask(ctx context.Context, taskID string) (*domain.Task, error) {
	if taskID == "" {
		return nil, fmt.Errorf("task_id is required")
	}

	// Validar que sea un UUID válido
	if _, err := uuid.FromString(taskID); err != nil {
		return nil, fmt.Errorf("invalid task_id format: %w", err)
	}

	return s.taskRepo.GetTask(ctx, taskID)
}

func (s *TaskService) UpdateTask(ctx context.Context, taskID string, title, description *string, completed *bool) (*domain.Task, error) {
	if taskID == "" {
		return nil, fmt.Errorf("task_id is required")
	}

	// Validar que sea un UUID válido
	if _, err := uuid.FromString(taskID); err != nil {
		return nil, fmt.Errorf("invalid task_id format: %w", err)
	}

	// Obtener la tarea existente
	existingTask, err := s.taskRepo.GetTask(ctx, taskID)
	if err != nil {
		return nil, err
	}

	// Actualizar solo los campos proporcionados
	if title != nil {
		existingTask.Title = *title
	}
	if description != nil {
		existingTask.Description = *description
	}
	if completed != nil {
		existingTask.Completed = *completed
	}

	return s.taskRepo.UpdateTask(ctx, existingTask)
}

func (s *TaskService) DeleteTask(ctx context.Context, taskID string) error {
	if taskID == "" {
		return fmt.Errorf("task_id is required")
	}

	// Validar que sea un UUID válido
	if _, err := uuid.FromString(taskID); err != nil {
		return fmt.Errorf("invalid task_id format: %w", err)
	}

	return s.taskRepo.DeleteTask(ctx, taskID)
}

func (s *TaskService) MarkTaskComplete(ctx context.Context, taskID string) (*domain.Task, error) {
	if taskID == "" {
		return nil, fmt.Errorf("task_id is required")
	}

	// Validar que sea un UUID válido
	if _, err := uuid.FromString(taskID); err != nil {
		return nil, fmt.Errorf("invalid task_id format: %w", err)
	}

	return s.taskRepo.MarkTaskComplete(ctx, taskID)
}

func (s *TaskService) ListTasksByUser(ctx context.Context, userID string) ([]*domain.Task, error) {
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	return s.taskRepo.ListTasksByUser(ctx, userID)
}

func (s *TaskService) ListAllTasks(ctx context.Context) ([]*domain.Task, error) {
	return s.taskRepo.ListAllTasks(ctx)
}
