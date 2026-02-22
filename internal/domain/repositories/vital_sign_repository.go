package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/kendall/chart-attack/internal/domain/entities"
)

type VitalSignRepository interface {
	Save(ctx context.Context, vs *entities.ValidatedVitalSign) error
	FindByPatientID(ctx context.Context, patientId uuid.UUID) ([]*entities.VitalSign, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entities.VitalSign, error)
	FindLatestByPatientID(ctx context.Context, patientId uuid.UUID) (*entities.VitalSign, error)
}
