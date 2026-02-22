package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/kendall/chart-attack/internal/domain/entities"
)

type EmergencyContactRepository interface {
	Save(ctx context.Context, ec *entities.ValidatedEmergencyContact) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByPatientID(ctx context.Context, patientId uuid.UUID) ([]*entities.EmergencyContact, error)
}
