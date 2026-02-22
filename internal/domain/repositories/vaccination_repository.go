package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/kendall/chart-attack/internal/domain/entities"
)

type VaccinationRepository interface {
	Save(ctx context.Context, vac *entities.ValidatedVaccinationRecord) error
	FindByPatientID(ctx context.Context, patientId uuid.UUID) ([]*entities.VaccinationRecord, error)
}
