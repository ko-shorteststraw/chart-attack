package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kendall/chart-attack/internal/domain/entities"
	"github.com/kendall/chart-attack/internal/domain/repositories"
	"github.com/kendall/chart-attack/internal/infrastructure/db/sqlc"
)

type sqlcTaskRepository struct {
	q *sqlc.Queries
}

func NewSqlcTaskRepository(db *sql.DB) repositories.TaskRepository {
	return &sqlcTaskRepository{q: sqlc.New(db)}
}

func (r *sqlcTaskRepository) Save(ctx context.Context, vt *entities.ValidatedTask) error {
	t := vt.Task()
	var completedAt sql.NullString
	if t.CompletedAt != nil {
		completedAt = sql.NullString{String: t.CompletedAt.Format(time.RFC3339), Valid: true}
	}
	var completedBy sql.NullString
	if t.CompletedBy != nil {
		completedBy = sql.NullString{String: t.CompletedBy.String(), Valid: true}
	}
	var recurInterval sql.NullString
	if t.RecurInterval != nil {
		recurInterval = sql.NullString{String: *t.RecurInterval, Valid: true}
	}
	return r.q.CreateTask(ctx, sqlc.CreateTaskParams{
		ID:            t.Id.String(),
		CreatedAt:     t.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     t.UpdatedAt.Format(time.RFC3339),
		PatientID:     t.PatientId.String(),
		AssignedTo:    t.AssignedTo.String(),
		Title:         t.Title,
		Category:      t.Category,
		DueAt:         t.DueAt.Format(time.RFC3339),
		CompletedAt:   completedAt,
		CompletedBy:   completedBy,
		Priority:      t.Priority,
		Recurring:     boolToInt64(t.Recurring),
		RecurInterval: recurInterval,
		Notes:         t.Notes,
	})
}

func (r *sqlcTaskRepository) Update(ctx context.Context, vt *entities.ValidatedTask) error {
	t := vt.Task()
	var completedAt sql.NullString
	if t.CompletedAt != nil {
		completedAt = sql.NullString{String: t.CompletedAt.Format(time.RFC3339), Valid: true}
	}
	var completedBy sql.NullString
	if t.CompletedBy != nil {
		completedBy = sql.NullString{String: t.CompletedBy.String(), Valid: true}
	}
	return r.q.UpdateTask(ctx, t.UpdatedAt.Format(time.RFC3339), completedAt, completedBy, t.Notes, t.Id.String())
}

func (r *sqlcTaskRepository) FindByPatientID(ctx context.Context, patientId uuid.UUID) ([]*entities.Task, error) {
	rows, err := r.q.GetTasksByPatientID(ctx, patientId.String())
	if err != nil {
		return nil, fmt.Errorf("finding tasks: %w", err)
	}
	tasks := make([]*entities.Task, 0, len(rows))
	for _, row := range rows {
		t, err := mapSqlcTaskToEntity(row)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *sqlcTaskRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Task, error) {
	row, err := r.q.GetTaskByID(ctx, id.String())
	if err != nil {
		return nil, fmt.Errorf("finding task by ID: %w", err)
	}
	return mapSqlcTaskToEntity(row)
}

func mapSqlcTaskToEntity(row sqlc.Task) (*entities.Task, error) {
	id, err := uuid.Parse(row.ID)
	if err != nil {
		return nil, fmt.Errorf("parsing task ID: %w", err)
	}
	createdAt, err := time.Parse(time.RFC3339, row.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("parsing created_at: %w", err)
	}
	updatedAt, err := time.Parse(time.RFC3339, row.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("parsing updated_at: %w", err)
	}
	patientId, err := uuid.Parse(row.PatientID)
	if err != nil {
		return nil, fmt.Errorf("parsing patient_id: %w", err)
	}
	assignedTo, err := uuid.Parse(row.AssignedTo)
	if err != nil {
		return nil, fmt.Errorf("parsing assigned_to: %w", err)
	}
	dueAt, err := time.Parse(time.RFC3339, row.DueAt)
	if err != nil {
		return nil, fmt.Errorf("parsing due_at: %w", err)
	}

	var completedAt *time.Time
	if row.CompletedAt.Valid {
		t, err := time.Parse(time.RFC3339, row.CompletedAt.String)
		if err == nil {
			completedAt = &t
		}
	}
	var completedBy *uuid.UUID
	if row.CompletedBy.Valid {
		uid, err := uuid.Parse(row.CompletedBy.String)
		if err == nil {
			completedBy = &uid
		}
	}
	var recurInterval *string
	if row.RecurInterval.Valid {
		recurInterval = &row.RecurInterval.String
	}

	return &entities.Task{
		Id:            id,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
		PatientId:     patientId,
		AssignedTo:    assignedTo,
		Title:         row.Title,
		Category:      row.Category,
		DueAt:         dueAt,
		CompletedAt:   completedAt,
		CompletedBy:   completedBy,
		Priority:      row.Priority,
		Recurring:     row.Recurring == 1,
		RecurInterval: recurInterval,
		Notes:         row.Notes,
	}, nil
}
