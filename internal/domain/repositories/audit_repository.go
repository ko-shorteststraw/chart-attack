package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/kendall/chart-attack/internal/domain/entities"
)

type AuditRepository interface {
	Save(ctx context.Context, entry *entities.AuditEntry) error
	FindByPatientID(ctx context.Context, patientId uuid.UUID) ([]*entities.AuditEntry, error)
	FindByEntity(ctx context.Context, entityType string, entityId uuid.UUID) ([]*entities.AuditEntry, error)
}
