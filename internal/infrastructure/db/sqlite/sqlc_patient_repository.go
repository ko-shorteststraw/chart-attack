package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kendall/chart-attack/internal/domain/entities"
	"github.com/kendall/chart-attack/internal/domain/repositories"
	"github.com/kendall/chart-attack/internal/infrastructure/db/sqlc"
)

type sqlcPatientRepository struct {
	q *sqlc.Queries
}

func NewSqlcPatientRepository(db *sql.DB) repositories.PatientRepository {
	return &sqlcPatientRepository{q: sqlc.New(db)}
}

func (r *sqlcPatientRepository) Save(ctx context.Context, vp *entities.ValidatedPatient) error {
	p := vp.Patient()
	allergiesJSON, _ := json.Marshal(p.Allergies)
	diagnosesJSON, _ := json.Marshal(p.Diagnoses)

	fallRisk := int64(0)
	if p.FallRisk {
		fallRisk = 1
	}

	var dischargeDate sql.NullString
	if p.DischargeDate != nil {
		dischargeDate = sql.NullString{String: p.DischargeDate.Format(time.RFC3339), Valid: true}
	}

	var assignedNurseID sql.NullString
	if p.AssignedNurseId != nil {
		assignedNurseID = sql.NullString{String: p.AssignedNurseId.String(), Valid: true}
	}

	return r.q.CreatePatient(ctx, sqlc.CreatePatientParams{
		ID:              p.Id.String(),
		CreatedAt:       p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       p.UpdatedAt.Format(time.RFC3339),
		Mrn:             p.MRN,
		FirstName:       p.FirstName,
		LastName:        p.LastName,
		DateOfBirth:     p.DateOfBirth.Format(time.RFC3339),
		RoomBed:         p.RoomBed,
		Allergies:       string(allergiesJSON),
		Diagnoses:       string(diagnosesJSON),
		CodeStatus:      p.CodeStatus,
		FallRisk:        fallRisk,
		IsolationType:   p.IsolationType,
		AdmitDate:       p.AdmitDate.Format(time.RFC3339),
		DischargeDate:   dischargeDate,
		AssignedNurseID: assignedNurseID,
	})
}

func (r *sqlcPatientRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Patient, error) {
	row, err := r.q.GetPatientByID(ctx, id.String())
	if err != nil {
		return nil, fmt.Errorf("finding patient by ID: %w", err)
	}
	return mapSqlcPatientToEntity(row)
}

func (r *sqlcPatientRepository) FindAdmitted(ctx context.Context) ([]*entities.Patient, error) {
	rows, err := r.q.GetAdmittedPatients(ctx)
	if err != nil {
		return nil, fmt.Errorf("finding admitted patients: %w", err)
	}
	patients := make([]*entities.Patient, 0, len(rows))
	for _, row := range rows {
		p, err := mapSqlcPatientToEntity(row)
		if err != nil {
			return nil, err
		}
		patients = append(patients, p)
	}
	return patients, nil
}

func (r *sqlcPatientRepository) Update(ctx context.Context, vp *entities.ValidatedPatient) error {
	p := vp.Patient()
	allergiesJSON, _ := json.Marshal(p.Allergies)
	diagnosesJSON, _ := json.Marshal(p.Diagnoses)

	fallRisk := int64(0)
	if p.FallRisk {
		fallRisk = 1
	}

	var assignedNurseID sql.NullString
	if p.AssignedNurseId != nil {
		assignedNurseID = sql.NullString{String: p.AssignedNurseId.String(), Valid: true}
	}

	return r.q.UpdatePatient(ctx, sqlc.UpdatePatientParams{
		UpdatedAt:       p.UpdatedAt.Format(time.RFC3339),
		FirstName:       p.FirstName,
		LastName:        p.LastName,
		RoomBed:         p.RoomBed,
		Allergies:       string(allergiesJSON),
		Diagnoses:       string(diagnosesJSON),
		CodeStatus:      p.CodeStatus,
		FallRisk:        fallRisk,
		IsolationType:   p.IsolationType,
		AssignedNurseID: assignedNurseID,
		ID:              p.Id.String(),
	})
}

func (r *sqlcPatientRepository) Discharge(ctx context.Context, p *entities.Patient) error {
	var dischargeDate sql.NullString
	if p.DischargeDate != nil {
		dischargeDate = sql.NullString{String: p.DischargeDate.Format(time.RFC3339), Valid: true}
	}
	return r.q.DischargePatient(ctx,
		p.UpdatedAt.Format(time.RFC3339),
		dischargeDate,
		p.Id.String(),
	)
}

func mapSqlcPatientToEntity(row sqlc.Patient) (*entities.Patient, error) {
	id, err := uuid.Parse(row.ID)
	if err != nil {
		return nil, fmt.Errorf("parsing patient ID: %w", err)
	}
	createdAt, err := time.Parse(time.RFC3339, row.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("parsing created_at: %w", err)
	}
	updatedAt, err := time.Parse(time.RFC3339, row.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("parsing updated_at: %w", err)
	}
	dob, err := time.Parse(time.RFC3339, row.DateOfBirth)
	if err != nil {
		return nil, fmt.Errorf("parsing date_of_birth: %w", err)
	}
	admitDate, err := time.Parse(time.RFC3339, row.AdmitDate)
	if err != nil {
		return nil, fmt.Errorf("parsing admit_date: %w", err)
	}

	var allergies []string
	if err := json.Unmarshal([]byte(row.Allergies), &allergies); err != nil {
		allergies = []string{}
	}
	var diagnoses []string
	if err := json.Unmarshal([]byte(row.Diagnoses), &diagnoses); err != nil {
		diagnoses = []string{}
	}

	var dischargeDate *time.Time
	if row.DischargeDate.Valid {
		t, err := time.Parse(time.RFC3339, row.DischargeDate.String)
		if err == nil {
			dischargeDate = &t
		}
	}

	var assignedNurseId *uuid.UUID
	if row.AssignedNurseID.Valid {
		nid, err := uuid.Parse(row.AssignedNurseID.String)
		if err == nil {
			assignedNurseId = &nid
		}
	}

	return &entities.Patient{
		Id:              id,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
		MRN:             row.Mrn,
		FirstName:       row.FirstName,
		LastName:        row.LastName,
		DateOfBirth:     dob,
		RoomBed:         row.RoomBed,
		Allergies:       allergies,
		Diagnoses:       diagnoses,
		CodeStatus:      row.CodeStatus,
		FallRisk:        row.FallRisk == 1,
		IsolationType:   row.IsolationType,
		AdmitDate:       admitDate,
		DischargeDate:   dischargeDate,
		AssignedNurseId: assignedNurseId,
	}, nil
}
