package interfaces

import (
	"context"

	"github.com/kendall/chart-attack/internal/application/command"
	"github.com/kendall/chart-attack/internal/application/query"
)

type TaskService interface {
	CreateTask(ctx context.Context, cmd command.CreateTaskCommand) (*command.CreateTaskResult, error)
	CompleteTask(ctx context.Context, cmd command.CompleteTaskCommand) error
	GetPatientTasks(ctx context.Context, q query.GetPatientTasksQuery) (*query.GetPatientTasksResult, error)
}
