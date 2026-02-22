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

type sqlcHealthHistoryRepository struct {
	q *sqlc.Queries
}

func NewSqlcHealthHistoryRepository(db *sql.DB) repositories.HealthHistoryRepository {
	return &sqlcHealthHistoryRepository{q: sqlc.New(db)}
}

func (r *sqlcHealthHistoryRepository) Save(ctx context.Context, vh *entities.ValidatedHealthHistoryEntry) error {
	h := vh.HealthHistoryEntry()
	var dateOccurred sql.NullString
	if h.DateOccurred != nil {
		dateOccurred = sql.NullString{String: h.DateOccurred.Format(time.RFC3339), Valid: true}
	}
	return r.q.CreateHealthHistory(ctx, sqlc.CreateHealthHistoryParams{
		ID:           h.Id.String(),
		CreatedAt:    h.CreatedAt.Format(time.RFC3339),
		PatientID:    h.PatientId.String(),
		RecordedBy:   h.RecordedBy.String(),
		Condition:    h.Condition,
		DateOccurred: dateOccurred,
		Description:  h.Description,
		Status:       h.Status,
	})
}

func (r *sqlcHealthHistoryRepository) FindByPatientID(ctx context.Context, patientId uuid.UUID) ([]*entities.HealthHistoryEntry, error) {
	rows, err := r.q.GetHealthHistoryByPatientID(ctx, patientId.String())
	if err != nil {
		return nil, fmt.Errorf("finding health history: %w", err)
	}
	entries := make([]*entities.HealthHistoryEntry, 0, len(rows))
	for _, row := range rows {
		e, err := mapSqlcHealthHistoryToEntity(row)
		if err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func mapSqlcHealthHistoryToEntity(row sqlc.HealthHistory) (*entities.HealthHistoryEntry, error) {
	id, err := uuid.Parse(row.ID)
	if err != nil {
		return nil, fmt.Errorf("parsing health history ID: %w", err)
	}
	createdAt, err := time.Parse(time.RFC3339, row.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("parsing created_at: %w", err)
	}
	patientId, err := uuid.Parse(row.PatientID)
	if err != nil {
		return nil, fmt.Errorf("parsing patient_id: %w", err)
	}
	recordedBy, err := uuid.Parse(row.RecordedBy)
	if err != nil {
		return nil, fmt.Errorf("parsing recorded_by: %w", err)
	}
	var dateOccurred *time.Time
	if row.DateOccurred.Valid {
		t, err := time.Parse(time.RFC3339, row.DateOccurred.String)
		if err == nil {
			dateOccurred = &t
		}
	}
	return &entities.HealthHistoryEntry{
		Id:           id,
		CreatedAt:    createdAt,
		PatientId:    patientId,
		RecordedBy:   recordedBy,
		Condition:    row.Condition,
		DateOccurred: dateOccurred,
		Description:  row.Description,
		Status:       row.Status,
	}, nil
}
