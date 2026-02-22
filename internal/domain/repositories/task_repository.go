package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/kendall/chart-attack/internal/domain/entities"
)

type TaskRepository interface {
	Save(ctx context.Context, task *entities.ValidatedTask) error
	Update(ctx context.Context, task *entities.ValidatedTask) error
	FindByPatientID(ctx context.Context, patientId uuid.UUID) ([]*entities.Task, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entities.Task, error)
}
