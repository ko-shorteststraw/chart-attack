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

type sqlcVitalSignRepository struct {
	q *sqlc.Queries
}

func NewSqlcVitalSignRepository(db *sql.DB) repositories.VitalSignRepository {
	return &sqlcVitalSignRepository{q: sqlc.New(db)}
}

func (r *sqlcVitalSignRepository) Save(ctx context.Context, vvs *entities.ValidatedVitalSign) error {
	vs := vvs.VitalSign()
	return r.q.CreateVitalSign(ctx, sqlc.CreateVitalSignParams{
		ID:             vs.Id.String(),
		CreatedAt:      vs.CreatedAt.Format(time.RFC3339),
		PatientID:      vs.PatientId.String(),
		RecordedBy:     vs.RecordedBy.String(),
		SystolicBp:     nullInt64FromIntPtr(vs.SystolicBP),
		DiastolicBp:    nullInt64FromIntPtr(vs.DiastolicBP),
		HeartRate:      nullInt64FromIntPtr(vs.HeartRate),
		Temperature:    nullFloat64FromPtr(vs.Temperature),
		TempRoute:      nullStringFromString(vs.TempRoute),
		OxygenSat:      nullInt64FromIntPtr(vs.OxygenSat),
		Respirations:   nullInt64FromIntPtr(vs.Respirations),
		PainLevel:      nullInt64FromIntPtr(vs.PainLevel),
		SupplementalO2: boolToInt64(vs.SupplementalO2),
		O2FlowRate:     nullFloat64FromPtr(vs.O2FlowRate),
		Position:       nullStringFromString(vs.Position),
		Notes:          vs.Notes,
	})
}

func (r *sqlcVitalSignRepository) FindByPatientID(ctx context.Context, patientId uuid.UUID) ([]*entities.VitalSign, error) {
	rows, err := r.q.GetVitalSignsByPatientID(ctx, patientId.String())
	if err != nil {
		return nil, fmt.Errorf("finding vital signs: %w", err)
	}
	vitals := make([]*entities.VitalSign, 0, len(rows))
	for _, row := range rows {
		vs, err := mapSqlcVitalSignToEntity(row)
		if err != nil {
			return nil, err
		}
		vitals = append(vitals, vs)
	}
	return vitals, nil
}

func (r *sqlcVitalSignRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.VitalSign, error) {
	row, err := r.q.GetVitalSignByID(ctx, id.String())
	if err != nil {
		return nil, fmt.Errorf("finding vital sign by ID: %w", err)
	}
	return mapSqlcVitalSignToEntity(row)
}

func (r *sqlcVitalSignRepository) FindLatestByPatientID(ctx context.Context, patientId uuid.UUID) (*entities.VitalSign, error) {
	row, err := r.q.GetLatestVitalSign(ctx, patientId.String())
	if err != nil {
		return nil, fmt.Errorf("finding latest vital sign: %w", err)
	}
	return mapSqlcVitalSignToEntity(row)
}

func mapSqlcVitalSignToEntity(row sqlc.VitalSign) (*entities.VitalSign, error) {
	id, err := uuid.Parse(row.ID)
	if err != nil {
		return nil, fmt.Errorf("parsing vital sign ID: %w", err)
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

	return &entities.VitalSign{
		Id:             id,
		CreatedAt:      createdAt,
		PatientId:      patientId,
		RecordedBy:     recordedBy,
		SystolicBP:     intPtrFromNullInt64(row.SystolicBp),
		DiastolicBP:    intPtrFromNullInt64(row.DiastolicBp),
		HeartRate:      intPtrFromNullInt64(row.HeartRate),
		Temperature:    float64PtrFromNull(row.Temperature),
		TempRoute:      row.TempRoute.String,
		OxygenSat:      intPtrFromNullInt64(row.OxygenSat),
		Respirations:   intPtrFromNullInt64(row.Respirations),
		PainLevel:      intPtrFromNullInt64(row.PainLevel),
		SupplementalO2: row.SupplementalO2 == 1,
		O2FlowRate:     float64PtrFromNull(row.O2FlowRate),
		Position:       row.Position.String,
		Notes:          row.Notes,
	}, nil
}

