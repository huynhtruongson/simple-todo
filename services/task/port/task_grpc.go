package task_port

import (
	"context"

	"github.com/huynhtruongson/simple-todo/common"
	"github.com/huynhtruongson/simple-todo/field"
	"github.com/huynhtruongson/simple-todo/interceptor"
	"github.com/huynhtruongson/simple-todo/pb"
	task_entity "github.com/huynhtruongson/simple-todo/services/task/entity"
	"github.com/huynhtruongson/simple-todo/token"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TaskGRPCService struct {
	pb.UnimplementedTaskServiceServer
	TaskService TaskService
}

func NewTaskGRPCService(taskService TaskService) *TaskGRPCService {
	return &TaskGRPCService{
		TaskService: taskService,
	}
}

func (sv *TaskGRPCService) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	payload, ok := ctx.Value(interceptor.AuthTokenPayload).(token.TokenPayload)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "Fail to parse token payload")
	}
	task := toTask(req)
	task.UserID = field.NewInt(payload.UserID)

	taskID, err := sv.TaskService.CreateTask(ctx, task)
	if err != nil {
		return nil, common.MapAppErrorToGRPCError(err, "Create task error")
	}
	return &pb.CreateTaskResponse{
		Data: int64(taskID),
	}, nil
}

func toTask(task *pb.CreateTaskRequest) task_entity.Task {
	description := field.NewNullString()
	if task.Description != nil {
		description = field.NewString(*task.Description)
	}
	return task_entity.Task{
		Title:       field.NewString(task.GetTitle()),
		Status:      task_entity.TaskStatus(task.GetStatus()),
		Description: description,
	}
}
