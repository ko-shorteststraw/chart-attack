package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/kendall/chart-attack/internal/domain/entities"
)

type MedicationRepository interface {
	Save(ctx context.Context, med *entities.ValidatedMedication) error
	FindByID(ctx context.Context, id uuid.UUID) (*entities.Medication, error)
	FindAll(ctx context.Context) ([]*entities.Medication, error)
}
