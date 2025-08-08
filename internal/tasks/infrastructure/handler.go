package infrastructure

import (
	"context"

	"github.com/Mayer-04/grpc-task-manager-go/internal/tasks/application"
	"github.com/Mayer-04/grpc-task-manager-go/internal/tasks/domain"
	"github.com/Mayer-04/grpc-task-manager-go/pkg/taskpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TaskHandler struct {
	taskpb.UnimplementedTaskServiceServer
	taskService *application.TaskService
}

func NewTaskHandler(taskService *application.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

func (h *TaskHandler) CreateTask(ctx context.Context, req *taskpb.CreateTaskRequest) (*taskpb.CreateTaskResponse, error) {
	completed := false
	if req.Completed != nil {
		completed = *req.Completed
	}

	description := ""
	if req.Description != nil {
		description = *req.Description
	}

	task, err := h.taskService.CreateTask(ctx, req.UserId, req.Title, description, completed)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create task: %v", err)
	}

	return &taskpb.CreateTaskResponse{
		Task: h.domainTaskToProto(task),
	}, nil
}

func (h *TaskHandler) GetTask(ctx context.Context, req *taskpb.GetTaskRequest) (*taskpb.GetTaskResponse, error) {
	task, err := h.taskService.GetTask(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "task not found: %v", err)
	}

	return &taskpb.GetTaskResponse{
		Task: h.domainTaskToProto(task),
	}, nil
}

func (h *TaskHandler) UpdateTask(ctx context.Context, req *taskpb.UpdateTaskRequest) (*taskpb.UpdateTaskResponse, error) {
	task, err := h.taskService.UpdateTask(ctx, req.Id, req.Title, req.Description, req.Completed)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update task: %v", err)
	}

	return &taskpb.UpdateTaskResponse{
		Task: h.domainTaskToProto(task),
	}, nil
}

func (h *TaskHandler) DeleteTask(ctx context.Context, req *taskpb.DeleteTaskRequest) (*taskpb.DeleteTaskResponse, error) {
	err := h.taskService.DeleteTask(ctx, req.Id)
	if err != nil {
		return &taskpb.DeleteTaskResponse{
			Success: false,
			Message: err.Error(),
		}, status.Errorf(codes.Internal, "failed to delete task: %v", err)
	}

	return &taskpb.DeleteTaskResponse{
		Success: true,
		Message: "Task deleted successfully",
	}, nil
}

func (h *TaskHandler) MarkTaskComplete(ctx context.Context, req *taskpb.MarkTaskCompleteRequest) (*taskpb.MarkTaskCompleteResponse, error) {
	task, err := h.taskService.MarkTaskComplete(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to mark task complete: %v", err)
	}

	return &taskpb.MarkTaskCompleteResponse{
		Task: h.domainTaskToProto(task),
	}, nil
}

func (h *TaskHandler) ListTasksByUser(ctx context.Context, req *taskpb.ListTasksByUserRequest) (*taskpb.ListTasksResponse, error) {
	tasks, err := h.taskService.ListTasksByUser(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list tasks: %v", err)
	}

	var protoTasks []*taskpb.Task
	for _, task := range tasks {
		protoTasks = append(protoTasks, h.domainTaskToProto(task))
	}

	return &taskpb.ListTasksResponse{
		Tasks: protoTasks,
	}, nil
}

func (h *TaskHandler) ListAllTasks(ctx context.Context, req *taskpb.ListAllTasksRequest) (*taskpb.ListTasksResponse, error) {
	tasks, err := h.taskService.ListAllTasks(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list all tasks: %v", err)
	}

	var protoTasks []*taskpb.Task
	for _, task := range tasks {
		protoTasks = append(protoTasks, h.domainTaskToProto(task))
	}

	return &taskpb.ListTasksResponse{
		Tasks: protoTasks,
	}, nil
}

// domainTaskToProto converts a domain.Task to a taskpb.Task.
func (h *TaskHandler) domainTaskToProto(task *domain.Task) *taskpb.Task {
	return &taskpb.Task{
		Id:        task.ID.String(),
		UserId:    task.UserID,
		Title:     task.Title,
		Completed: task.Completed,
		CreatedAt: timestamppb.New(task.CreatedAt),
		UpdatedAt: timestamppb.New(task.UpdatedAt),
	}
}
