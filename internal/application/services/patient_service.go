package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kendall/chart-attack/internal/application/command"
	appinterfaces "github.com/kendall/chart-attack/internal/application/interfaces"
	"github.com/kendall/chart-attack/internal/application/mapper"
	"github.com/kendall/chart-attack/internal/application/query"
	"github.com/kendall/chart-attack/internal/domain/entities"
	"github.com/kendall/chart-attack/internal/domain/repositories"
)

type patientService struct {
	patientRepo          repositories.PatientRepository
	emergencyContactRepo repositories.EmergencyContactRepository
	auditRepo            repositories.AuditRepository
	idempotencyRepo      repositories.IdempotencyRepository
}

func NewPatientService(
	patientRepo repositories.PatientRepository,
	emergencyContactRepo repositories.EmergencyContactRepository,
	auditRepo repositories.AuditRepository,
	idempotencyRepo repositories.IdempotencyRepository,
) appinterfaces.PatientService {
	return &patientService{
		patientRepo:          patientRepo,
		emergencyContactRepo: emergencyContactRepo,
		auditRepo:            auditRepo,
		idempotencyRepo:      idempotencyRepo,
	}
}

func (s *patientService) AdmitPatient(ctx context.Context, cmd command.AdmitPatientCommand) (*command.AdmitPatientResult, error) {
	// Check idempotency
	if cmd.IdempotencyKey != "" {
		existing, err := s.idempotencyRepo.FindByKey(ctx, cmd.IdempotencyKey)
		if err != nil {
			return nil, fmt.Errorf("checking idempotency: %w", err)
		}
		if existing != nil {
			var result command.AdmitPatientResult
			if err := json.Unmarshal([]byte(existing.Response), &result); err == nil {
				return &result, nil
			}
		}
	}

	patient, err := entities.NewPatient(
		cmd.MRN, cmd.FirstName, cmd.LastName, cmd.DateOfBirth,
		cmd.RoomBed, cmd.Allergies, cmd.Diagnoses,
		cmd.CodeStatus, cmd.IsolationType, cmd.FallRisk,
	)
	if err != nil {
		return nil, fmt.Errorf("creating patient: %w", err)
	}

	validated, err := entities.NewValidatedPatient(patient)
	if err != nil {
		return nil, fmt.Errorf("validating patient: %w", err)
	}

	if err := s.patientRepo.Save(ctx, validated); err != nil {
		return nil, fmt.Errorf("saving patient: %w", err)
	}

	// Audit trail
	audit := entities.NewAuditEntry(cmd.AdmittedBy, &patient.Id, "CREATE", "Patient", patient.Id, "{}", "", "")
	if err := s.auditRepo.Save(ctx, audit); err != nil {
		return nil, fmt.Errorf("saving audit entry: %w", err)
	}

	result := &command.AdmitPatientResult{PatientId: patient.Id.String()}

	// Save idempotency record
	if cmd.IdempotencyKey != "" {
		responseJSON, _ := json.Marshal(result)
		record := entities.NewIdempotencyRecord(cmd.IdempotencyKey, string(responseJSON))
		_ = s.idempotencyRepo.Save(ctx, record)
	}

	return result, nil
}

func (s *patientService) GetPatientById(ctx context.Context, q query.GetPatientByIdQuery) (*query.GetPatientByIdResult, error) {
	patient, err := s.patientRepo.FindByID(ctx, q.PatientId)
	if err != nil {
		return nil, fmt.Errorf("finding patient: %w", err)
	}
	return &query.GetPatientByIdResult{
		Patient: mapper.ToPatientResult(patient),
	}, nil
}

func (s *patientService) GetPatientCensus(ctx context.Context, q query.GetPatientCensusQuery) (*query.GetPatientCensusResult, error) {
	patients, err := s.patientRepo.FindAdmitted(ctx)
	if err != nil {
		return nil, fmt.Errorf("finding admitted patients: %w", err)
	}
	return &query.GetPatientCensusResult{
		Patients: mapper.ToPatientResults(patients),
	}, nil
}

func (s *patientService) UpdateDiagnoses(ctx context.Context, cmd command.UpdateDiagnosesCommand) error {
	patient, err := s.patientRepo.FindByID(ctx, cmd.PatientId)
	if err != nil {
		return fmt.Errorf("finding patient: %w", err)
	}

	switch cmd.Action {
	case "add":
		patient.AddDiagnosis(cmd.Diagnosis)
	case "remove":
		patient.RemoveDiagnosis(cmd.Diagnosis)
	default:
		return fmt.Errorf("invalid action: %s", cmd.Action)
	}

	validated, err := entities.NewValidatedPatient(patient)
	if err != nil {
		return fmt.Errorf("validating patient: %w", err)
	}

	if err := s.patientRepo.Update(ctx, validated); err != nil {
		return fmt.Errorf("updating patient: %w", err)
	}

	audit := entities.NewAuditEntry(cmd.UpdatedBy, &cmd.PatientId, "UPDATE", "Patient", cmd.PatientId, fmt.Sprintf(`{"diagnoses":"%s %s"}`, cmd.Action, cmd.Diagnosis), "", "")
	_ = s.auditRepo.Save(ctx, audit)

	return nil
}

func (s *patientService) AddEmergencyContact(ctx context.Context, cmd command.AddEmergencyContactCommand) (*command.AddEmergencyContactResult, error) {
	ec, err := entities.NewEmergencyContact(cmd.PatientId, cmd.Name, cmd.Relationship, cmd.Phone, cmd.Email, cmd.IsPrimary)
	if err != nil {
		return nil, fmt.Errorf("creating emergency contact: %w", err)
	}

	validated, err := entities.NewValidatedEmergencyContact(ec)
	if err != nil {
		return nil, fmt.Errorf("validating emergency contact: %w", err)
	}

	if err := s.emergencyContactRepo.Save(ctx, validated); err != nil {
		return nil, fmt.Errorf("saving emergency contact: %w", err)
	}

	audit := entities.NewAuditEntry(cmd.AddedBy, &cmd.PatientId, "CREATE", "EmergencyContact", ec.Id, "{}", "", "")
	_ = s.auditRepo.Save(ctx, audit)

	return &command.AddEmergencyContactResult{ContactId: ec.Id.String()}, nil
}

func (s *patientService) RemoveEmergencyContact(ctx context.Context, cmd command.RemoveEmergencyContactCommand) error {
	if err := s.emergencyContactRepo.Delete(ctx, cmd.ContactId); err != nil {
		return fmt.Errorf("deleting emergency contact: %w", err)
	}

	audit := entities.NewAuditEntry(cmd.RemovedBy, nil, "DELETE", "EmergencyContact", cmd.ContactId, "{}", "", "")
	_ = s.auditRepo.Save(ctx, audit)

	return nil
}

func (s *patientService) GetPatientEmergencyContacts(ctx context.Context, q query.GetPatientEmergencyContactsQuery) (*query.GetPatientEmergencyContactsResult, error) {
	contacts, err := s.emergencyContactRepo.FindByPatientID(ctx, q.PatientId)
	if err != nil {
		return nil, fmt.Errorf("finding emergency contacts: %w", err)
	}
	return &query.GetPatientEmergencyContactsResult{
		Contacts: mapper.ToEmergencyContactResults(contacts),
	}, nil
}
