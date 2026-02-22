package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/kendall/chart-attack/internal/domain/entities"
)

type MARRepository interface {
	Save(ctx context.Context, entry *entities.ValidatedMAREntry) error
	Update(ctx context.Context, entry *entities.ValidatedMAREntry) error
	FindByPatientID(ctx context.Context, patientId uuid.UUID) ([]*entities.MAREntry, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entities.MAREntry, error)
}
