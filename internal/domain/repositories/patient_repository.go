package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/kendall/chart-attack/internal/domain/entities"
)

type PatientRepository interface {
	Save(ctx context.Context, patient *entities.ValidatedPatient) error
	FindByID(ctx context.Context, id uuid.UUID) (*entities.Patient, error)
	FindAdmitted(ctx context.Context) ([]*entities.Patient, error)
	Update(ctx context.Context, patient *entities.ValidatedPatient) error
	Discharge(ctx context.Context, patient *entities.Patient) error
}
