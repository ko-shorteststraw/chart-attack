package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/kendall/chart-attack/internal/domain/entities"
)

type HealthHistoryRepository interface {
	Save(ctx context.Context, entry *entities.ValidatedHealthHistoryEntry) error
	FindByPatientID(ctx context.Context, patientId uuid.UUID) ([]*entities.HealthHistoryEntry, error)
}
