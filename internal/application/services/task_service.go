package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kendall/chart-attack/internal/application/command"
	appinterfaces "github.com/kendall/chart-attack/internal/application/interfaces"
	"github.com/kendall/chart-attack/internal/application/mapper"
	"github.com/kendall/chart-attack/internal/application/query"
	"github.com/kendall/chart-attack/internal/domain/entities"
	"github.com/kendall/chart-attack/internal/domain/repositories"
)

type taskService struct {
	taskRepo        repositories.TaskRepository
	auditRepo       repositories.AuditRepository
	idempotencyRepo repositories.IdempotencyRepository
}

func NewTaskService(
	taskRepo repositories.TaskRepository,
	auditRepo repositories.AuditRepository,
	idempotencyRepo repositories.IdempotencyRepository,
) appinterfaces.TaskService {
	return &taskService{
		taskRepo:        taskRepo,
		auditRepo:       auditRepo,
		idempotencyRepo: idempotencyRepo,
	}
}

func (s *taskService) CreateTask(ctx context.Context, cmd command.CreateTaskCommand) (*command.CreateTaskResult, error) {
	if cmd.IdempotencyKey != "" {
		existing, err := s.idempotencyRepo.FindByKey(ctx, cmd.IdempotencyKey)
		if err != nil {
			return nil, fmt.Errorf("checking idempotency: %w", err)
		}
		if existing != nil {
			var result command.CreateTaskResult
			if err := json.Unmarshal([]byte(existing.Response), &result); err == nil {
				return &result, nil
			}
		}
	}

	task, err := entities.NewTask(cmd.PatientId, cmd.AssignedTo, cmd.Title, cmd.Category, cmd.DueAt, cmd.Priority)
	if err != nil {
		return nil, fmt.Errorf("creating task: %w", err)
	}
	task.SetNotes(cmd.Notes)

	validated, err := entities.NewValidatedTask(task)
	if err != nil {
		return nil, fmt.Errorf("validating task: %w", err)
	}

	if err := s.taskRepo.Save(ctx, validated); err != nil {
		return nil, fmt.Errorf("saving task: %w", err)
	}

	audit := entities.NewAuditEntry(cmd.CreatedBy, &cmd.PatientId, "CREATE", "Task", task.Id, "{}", "", "")
	_ = s.auditRepo.Save(ctx, audit)

	result := &command.CreateTaskResult{TaskId: task.Id.String()}

	if cmd.IdempotencyKey != "" {
		responseJSON, _ := json.Marshal(result)
		record := entities.NewIdempotencyRecord(cmd.IdempotencyKey, string(responseJSON))
		_ = s.idempotencyRepo.Save(ctx, record)
	}

	return result, nil
}

func (s *taskService) CompleteTask(ctx context.Context, cmd command.CompleteTaskCommand) error {
	task, err := s.taskRepo.FindByID(ctx, cmd.TaskId)
	if err != nil {
		return fmt.Errorf("finding task: %w", err)
	}

	task.Complete(cmd.CompletedBy)

	validated, err := entities.NewValidatedTask(task)
	if err != nil {
		return fmt.Errorf("validating task: %w", err)
	}

	if err := s.taskRepo.Update(ctx, validated); err != nil {
		return fmt.Errorf("updating task: %w", err)
	}

	audit := entities.NewAuditEntry(cmd.CompletedBy, &task.PatientId, "UPDATE", "Task", task.Id, `{"status":"completed"}`, "", "")
	_ = s.auditRepo.Save(ctx, audit)

	return nil
}

func (s *taskService) GetPatientTasks(ctx context.Context, q query.GetPatientTasksQuery) (*query.GetPatientTasksResult, error) {
	tasks, err := s.taskRepo.FindByPatientID(ctx, q.PatientId)
	if err != nil {
		return nil, fmt.Errorf("finding tasks: %w", err)
	}
	return &query.GetPatientTasksResult{
		Tasks: mapper.ToTaskResults(tasks),
	}, nil
}
