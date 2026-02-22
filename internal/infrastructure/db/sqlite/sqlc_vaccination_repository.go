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

type sqlcVaccinationRepository struct {
	q *sqlc.Queries
}

func NewSqlcVaccinationRepository(db *sql.DB) repositories.VaccinationRepository {
	return &sqlcVaccinationRepository{q: sqlc.New(db)}
}

func (r *sqlcVaccinationRepository) Save(ctx context.Context, vv *entities.ValidatedVaccinationRecord) error {
	v := vv.VaccinationRecord()
	return r.q.CreateVaccinationRecord(ctx, sqlc.CreateVaccinationRecordParams{
		ID:               v.Id.String(),
		CreatedAt:        v.CreatedAt.Format(time.RFC3339),
		PatientID:        v.PatientId.String(),
		RecordedBy:       v.RecordedBy.String(),
		VaccineName:      v.VaccineName,
		DateAdministered: v.DateAdministered.Format(time.RFC3339),
		LotNumber:        v.LotNumber,
		Site:             v.Site,
		Notes:            v.Notes,
	})
}

func (r *sqlcVaccinationRepository) FindByPatientID(ctx context.Context, patientId uuid.UUID) ([]*entities.VaccinationRecord, error) {
	rows, err := r.q.GetVaccinationRecordsByPatientID(ctx, patientId.String())
	if err != nil {
		return nil, fmt.Errorf("finding vaccination records: %w", err)
	}
	records := make([]*entities.VaccinationRecord, 0, len(rows))
	for _, row := range rows {
		v, err := mapSqlcVaccinationToEntity(row)
		if err != nil {
			return nil, err
		}
		records = append(records, v)
	}
	return records, nil
}

func mapSqlcVaccinationToEntity(row sqlc.VaccinationRecord) (*entities.VaccinationRecord, error) {
	id, err := uuid.Parse(row.ID)
	if err != nil {
		return nil, fmt.Errorf("parsing vaccination ID: %w", err)
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
	dateAdministered, err := time.Parse(time.RFC3339, row.DateAdministered)
	if err != nil {
		return nil, fmt.Errorf("parsing date_administered: %w", err)
	}
	return &entities.VaccinationRecord{
		Id:               id,
		CreatedAt:        createdAt,
		PatientId:        patientId,
		RecordedBy:       recordedBy,
		VaccineName:      row.VaccineName,
		DateAdministered: dateAdministered,
		LotNumber:        row.LotNumber,
		Site:             row.Site,
		Notes:            row.Notes,
	}, nil
}
