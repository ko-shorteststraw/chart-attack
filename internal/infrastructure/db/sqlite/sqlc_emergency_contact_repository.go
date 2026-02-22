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

type sqlcEmergencyContactRepository struct {
	q *sqlc.Queries
}

func NewSqlcEmergencyContactRepository(db *sql.DB) repositories.EmergencyContactRepository {
	return &sqlcEmergencyContactRepository{q: sqlc.New(db)}
}

func (r *sqlcEmergencyContactRepository) Save(ctx context.Context, vec *entities.ValidatedEmergencyContact) error {
	ec := vec.EmergencyContact()
	return r.q.CreateEmergencyContact(ctx, sqlc.CreateEmergencyContactParams{
		ID:           ec.Id.String(),
		CreatedAt:    ec.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    ec.UpdatedAt.Format(time.RFC3339),
		PatientID:    ec.PatientId.String(),
		Name:         ec.Name,
		Relationship: ec.Relationship,
		Phone:        ec.Phone,
		Email:        ec.Email,
		IsPrimary:    boolToInt64(ec.IsPrimary),
	})
}

func (r *sqlcEmergencyContactRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteEmergencyContact(ctx, id.String())
}

func (r *sqlcEmergencyContactRepository) FindByPatientID(ctx context.Context, patientId uuid.UUID) ([]*entities.EmergencyContact, error) {
	rows, err := r.q.GetEmergencyContactsByPatientID(ctx, patientId.String())
	if err != nil {
		return nil, fmt.Errorf("finding emergency contacts: %w", err)
	}
	contacts := make([]*entities.EmergencyContact, 0, len(rows))
	for _, row := range rows {
		ec, err := mapSqlcEmergencyContactToEntity(row)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, ec)
	}
	return contacts, nil
}

func mapSqlcEmergencyContactToEntity(row sqlc.EmergencyContact) (*entities.EmergencyContact, error) {
	id, err := uuid.Parse(row.ID)
	if err != nil {
		return nil, fmt.Errorf("parsing emergency contact ID: %w", err)
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
	return &entities.EmergencyContact{
		Id:           id,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		PatientId:    patientId,
		Name:         row.Name,
		Relationship: row.Relationship,
		Phone:        row.Phone,
		Email:        row.Email,
		IsPrimary:    row.IsPrimary == 1,
	}, nil
}
