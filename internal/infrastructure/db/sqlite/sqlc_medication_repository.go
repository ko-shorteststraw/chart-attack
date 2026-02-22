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

type sqlcMedicationRepository struct {
	q *sqlc.Queries
}

func NewSqlcMedicationRepository(db *sql.DB) repositories.MedicationRepository {
	return &sqlcMedicationRepository{q: sqlc.New(db)}
}

func (r *sqlcMedicationRepository) Save(ctx context.Context, vm *entities.ValidatedMedication) error {
	m := vm.Medication()
	return r.q.CreateMedication(ctx, sqlc.CreateMedicationParams{
		ID:           m.Id.String(),
		CreatedAt:    m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    m.UpdatedAt.Format(time.RFC3339),
		Name:         m.Name,
		BrandName:    m.BrandName,
		DrugClass:    m.DrugClass,
		NdcCode:      m.NDCCode,
		DefaultDose:  m.DefaultDose,
		DefaultRoute: m.DefaultRoute,
		Frequency:    m.Frequency,
		HighAlert:    boolToInt64(m.HighAlert),
	})
}

func (r *sqlcMedicationRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Medication, error) {
	row, err := r.q.GetMedicationByID(ctx, id.String())
	if err != nil {
		return nil, fmt.Errorf("finding medication by ID: %w", err)
	}
	return mapSqlcMedicationToEntity(row)
}

func (r *sqlcMedicationRepository) FindAll(ctx context.Context) ([]*entities.Medication, error) {
	rows, err := r.q.GetAllMedications(ctx)
	if err != nil {
		return nil, fmt.Errorf("finding all medications: %w", err)
	}
	meds := make([]*entities.Medication, 0, len(rows))
	for _, row := range rows {
		m, err := mapSqlcMedicationToEntity(row)
		if err != nil {
			return nil, err
		}
		meds = append(meds, m)
	}
	return meds, nil
}

func mapSqlcMedicationToEntity(row sqlc.Medication) (*entities.Medication, error) {
	id, err := uuid.Parse(row.ID)
	if err != nil {
		return nil, fmt.Errorf("parsing medication ID: %w", err)
	}
	createdAt, err := time.Parse(time.RFC3339, row.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("parsing created_at: %w", err)
	}
	updatedAt, err := time.Parse(time.RFC3339, row.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("parsing updated_at: %w", err)
	}
	return &entities.Medication{
		Id:           id,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		Name:         row.Name,
		BrandName:    row.BrandName,
		DrugClass:    row.DrugClass,
		NDCCode:      row.NdcCode,
		DefaultDose:  row.DefaultDose,
		DefaultRoute: row.DefaultRoute,
		Frequency:    row.Frequency,
		HighAlert:    row.HighAlert == 1,
	}, nil
}
