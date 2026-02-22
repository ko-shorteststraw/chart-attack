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

type sqlcAuditRepository struct {
	q *sqlc.Queries
}

func NewSqlcAuditRepository(db *sql.DB) repositories.AuditRepository {
	return &sqlcAuditRepository{q: sqlc.New(db)}
}

func (r *sqlcAuditRepository) Save(ctx context.Context, entry *entities.AuditEntry) error {
	var patientID sql.NullString
	if entry.PatientId != nil {
		patientID = sql.NullString{String: entry.PatientId.String(), Valid: true}
	}
	return r.q.CreateAuditEntry(ctx, sqlc.CreateAuditEntryParams{
		ID:            entry.Id.String(),
		CreatedAt:     entry.CreatedAt.Format(time.RFC3339),
		UserID:        entry.UserId.String(),
		PatientID:     patientID,
		Action:        entry.Action,
		EntityType:    entry.EntityType,
		EntityID:      entry.EntityId.String(),
		FieldsChanged: entry.FieldsChanged,
		IPAddress:     entry.IPAddress,
		UserAgent:     entry.UserAgent,
	})
}

func (r *sqlcAuditRepository) FindByPatientID(ctx context.Context, patientId uuid.UUID) ([]*entities.AuditEntry, error) {
	rows, err := r.q.GetAuditEntriesByPatientID(ctx, sql.NullString{String: patientId.String(), Valid: true})
	if err != nil {
		return nil, fmt.Errorf("finding audit entries: %w", err)
	}
	return mapSqlcAuditEntries(rows)
}

func (r *sqlcAuditRepository) FindByEntity(ctx context.Context, entityType string, entityId uuid.UUID) ([]*entities.AuditEntry, error) {
	rows, err := r.q.GetAuditEntriesByEntity(ctx, entityType, entityId.String())
	if err != nil {
		return nil, fmt.Errorf("finding audit entries by entity: %w", err)
	}
	return mapSqlcAuditEntries(rows)
}

func mapSqlcAuditEntries(rows []sqlc.AuditEntry) ([]*entities.AuditEntry, error) {
	entries := make([]*entities.AuditEntry, 0, len(rows))
	for _, row := range rows {
		e, err := mapSqlcAuditEntryToEntity(row)
		if err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func mapSqlcAuditEntryToEntity(row sqlc.AuditEntry) (*entities.AuditEntry, error) {
	id, err := uuid.Parse(row.ID)
	if err != nil {
		return nil, fmt.Errorf("parsing audit entry ID: %w", err)
	}
	createdAt, err := time.Parse(time.RFC3339, row.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("parsing created_at: %w", err)
	}
	userId, err := uuid.Parse(row.UserID)
	if err != nil {
		return nil, fmt.Errorf("parsing user_id: %w", err)
	}
	entityId, err := uuid.Parse(row.EntityID)
	if err != nil {
		return nil, fmt.Errorf("parsing entity_id: %w", err)
	}

	var patientId *uuid.UUID
	if row.PatientID.Valid {
		pid, err := uuid.Parse(row.PatientID.String)
		if err == nil {
			patientId = &pid
		}
	}

	return &entities.AuditEntry{
		Id:            id,
		CreatedAt:     createdAt,
		UserId:        userId,
		PatientId:     patientId,
		Action:        row.Action,
		EntityType:    row.EntityType,
		EntityId:      entityId,
		FieldsChanged: row.FieldsChanged,
		IPAddress:     row.IPAddress,
		UserAgent:     row.UserAgent,
	}, nil
}
