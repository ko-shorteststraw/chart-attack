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

type sqlcMARRepository struct {
	q *sqlc.Queries
}

func NewSqlcMARRepository(db *sql.DB) repositories.MARRepository {
	return &sqlcMARRepository{q: sqlc.New(db)}
}

func (r *sqlcMARRepository) Save(ctx context.Context, vm *entities.ValidatedMAREntry) error {
	m := vm.MAREntry()
	var administeredAt sql.NullString
	if m.AdministeredAt != nil {
		administeredAt = sql.NullString{String: m.AdministeredAt.Format(time.RFC3339), Valid: true}
	}
	var administeredBy sql.NullString
	if m.AdministeredBy != nil {
		administeredBy = sql.NullString{String: m.AdministeredBy.String(), Valid: true}
	}

	return r.q.CreateMAREntry(ctx, sqlc.CreateMAREntryParams{
		ID:             m.Id.String(),
		CreatedAt:      m.CreatedAt.Format(time.RFC3339),
		PatientID:      m.PatientId.String(),
		MedicationID:   m.MedicationId.String(),
		ScheduledTime:  m.ScheduledTime.Format(time.RFC3339),
		AdministeredAt: administeredAt,
		AdministeredBy: administeredBy,
		Status:         m.Status,
		Dose:           m.Dose,
		Route:          m.Route,
		Site:           m.Site,
		HoldReason:     m.HoldReason,
		Notes:          m.Notes,
	})
}

func (r *sqlcMARRepository) Update(ctx context.Context, vm *entities.ValidatedMAREntry) error {
	m := vm.MAREntry()
	var administeredAt sql.NullString
	if m.AdministeredAt != nil {
		administeredAt = sql.NullString{String: m.AdministeredAt.Format(time.RFC3339), Valid: true}
	}
	var administeredBy sql.NullString
	if m.AdministeredBy != nil {
		administeredBy = sql.NullString{String: m.AdministeredBy.String(), Valid: true}
	}
	return r.q.UpdateMAREntry(ctx, administeredAt, administeredBy, m.Status, m.HoldReason, m.Notes, m.Id.String())
}

func (r *sqlcMARRepository) FindByPatientID(ctx context.Context, patientId uuid.UUID) ([]*entities.MAREntry, error) {
	rows, err := r.q.GetMAREntriesByPatientID(ctx, patientId.String())
	if err != nil {
		return nil, fmt.Errorf("finding MAR entries: %w", err)
	}
	entries := make([]*entities.MAREntry, 0, len(rows))
	for _, row := range rows {
		e, err := mapSqlcMAREntryToEntity(row)
		if err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func (r *sqlcMARRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.MAREntry, error) {
	row, err := r.q.GetMAREntryByID(ctx, id.String())
	if err != nil {
		return nil, fmt.Errorf("finding MAR entry by ID: %w", err)
	}
	return mapSqlcMAREntryToEntity(row)
}

func mapSqlcMAREntryToEntity(row sqlc.MAREntry) (*entities.MAREntry, error) {
	id, err := uuid.Parse(row.ID)
	if err != nil {
		return nil, fmt.Errorf("parsing MAR entry ID: %w", err)
	}
	createdAt, err := time.Parse(time.RFC3339, row.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("parsing created_at: %w", err)
	}
	patientId, err := uuid.Parse(row.PatientID)
	if err != nil {
		return nil, fmt.Errorf("parsing patient_id: %w", err)
	}
	medicationId, err := uuid.Parse(row.MedicationID)
	if err != nil {
		return nil, fmt.Errorf("parsing medication_id: %w", err)
	}
	scheduledTime, err := time.Parse(time.RFC3339, row.ScheduledTime)
	if err != nil {
		return nil, fmt.Errorf("parsing scheduled_time: %w", err)
	}

	var administeredAt *time.Time
	if row.AdministeredAt.Valid {
		t, err := time.Parse(time.RFC3339, row.AdministeredAt.String)
		if err == nil {
			administeredAt = &t
		}
	}
	var administeredBy *uuid.UUID
	if row.AdministeredBy.Valid {
		uid, err := uuid.Parse(row.AdministeredBy.String)
		if err == nil {
			administeredBy = &uid
		}
	}

	return &entities.MAREntry{
		Id:             id,
		CreatedAt:      createdAt,
		PatientId:      patientId,
		MedicationId:   medicationId,
		ScheduledTime:  scheduledTime,
		AdministeredAt: administeredAt,
		AdministeredBy: administeredBy,
		Status:         row.Status,
		Dose:           row.Dose,
		Route:          row.Route,
		Site:           row.Site,
		HoldReason:     row.HoldReason,
		Notes:          row.Notes,
	}, nil
}
