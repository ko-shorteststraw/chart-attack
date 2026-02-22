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

type medicationService struct {
	medicationRepo  repositories.MedicationRepository
	marRepo         repositories.MARRepository
	auditRepo       repositories.AuditRepository
	idempotencyRepo repositories.IdempotencyRepository
}

func NewMedicationService(
	medicationRepo repositories.MedicationRepository,
	marRepo repositories.MARRepository,
	auditRepo repositories.AuditRepository,
	idempotencyRepo repositories.IdempotencyRepository,
) appinterfaces.MedicationService {
	return &medicationService{
		medicationRepo:  medicationRepo,
		marRepo:         marRepo,
		auditRepo:       auditRepo,
		idempotencyRepo: idempotencyRepo,
	}
}

func (s *medicationService) AddMedication(ctx context.Context, cmd command.AddMedicationCommand) (*command.AddMedicationResult, error) {
	if cmd.IdempotencyKey != "" {
		existing, err := s.idempotencyRepo.FindByKey(ctx, cmd.IdempotencyKey)
		if err != nil {
			return nil, fmt.Errorf("checking idempotency: %w", err)
		}
		if existing != nil {
			var result command.AddMedicationResult
			if err := json.Unmarshal([]byte(existing.Response), &result); err == nil {
				return &result, nil
			}
		}
	}

	med, err := entities.NewMedication(cmd.Name, cmd.BrandName, cmd.DrugClass, cmd.DefaultDose, cmd.DefaultRoute, cmd.Frequency, cmd.HighAlert)
	if err != nil {
		return nil, fmt.Errorf("creating medication: %w", err)
	}

	validated, err := entities.NewValidatedMedication(med)
	if err != nil {
		return nil, fmt.Errorf("validating medication: %w", err)
	}

	if err := s.medicationRepo.Save(ctx, validated); err != nil {
		return nil, fmt.Errorf("saving medication: %w", err)
	}

	audit := entities.NewAuditEntry(cmd.AddedBy, nil, "CREATE", "Medication", med.Id, "{}", "", "")
	_ = s.auditRepo.Save(ctx, audit)

	result := &command.AddMedicationResult{MedicationId: med.Id.String()}

	if cmd.IdempotencyKey != "" {
		responseJSON, _ := json.Marshal(result)
		record := entities.NewIdempotencyRecord(cmd.IdempotencyKey, string(responseJSON))
		_ = s.idempotencyRepo.Save(ctx, record)
	}

	return result, nil
}

func (s *medicationService) AdministerMedication(ctx context.Context, cmd command.AdministerMedicationCommand) (*command.AdministerMedicationResult, error) {
	if cmd.IdempotencyKey != "" {
		existing, err := s.idempotencyRepo.FindByKey(ctx, cmd.IdempotencyKey)
		if err != nil {
			return nil, fmt.Errorf("checking idempotency: %w", err)
		}
		if existing != nil {
			var result command.AdministerMedicationResult
			if err := json.Unmarshal([]byte(existing.Response), &result); err == nil {
				return &result, nil
			}
		}
	}

	entry, err := entities.NewMAREntry(cmd.PatientId, cmd.MedicationId, cmd.ScheduledTime, cmd.Dose, cmd.Route)
	if err != nil {
		return nil, fmt.Errorf("creating MAR entry: %w", err)
	}

	if err := entry.Administer(cmd.AdministeredBy); err != nil {
		return nil, fmt.Errorf("administering medication: %w", err)
	}

	validated, err := entities.NewValidatedMAREntry(entry)
	if err != nil {
		return nil, fmt.Errorf("validating MAR entry: %w", err)
	}

	if err := s.marRepo.Save(ctx, validated); err != nil {
		return nil, fmt.Errorf("saving MAR entry: %w", err)
	}

	audit := entities.NewAuditEntry(cmd.AdministeredBy, &cmd.PatientId, "CREATE", "MAREntry", entry.Id, "{}", "", "")
	_ = s.auditRepo.Save(ctx, audit)

	result := &command.AdministerMedicationResult{MAREntryId: entry.Id.String()}

	if cmd.IdempotencyKey != "" {
		responseJSON, _ := json.Marshal(result)
		record := entities.NewIdempotencyRecord(cmd.IdempotencyKey, string(responseJSON))
		_ = s.idempotencyRepo.Save(ctx, record)
	}

	return result, nil
}

func (s *medicationService) GetPatientMAR(ctx context.Context, q query.GetPatientMARQuery) (*query.GetPatientMARResult, error) {
	entries, err := s.marRepo.FindByPatientID(ctx, q.PatientId)
	if err != nil {
		return nil, fmt.Errorf("finding MAR entries: %w", err)
	}

	// Build medication name map
	medNames := make(map[string]string)
	for _, e := range entries {
		midStr := e.MedicationId.String()
		if _, ok := medNames[midStr]; !ok {
			med, err := s.medicationRepo.FindByID(ctx, e.MedicationId)
			if err == nil {
				medNames[midStr] = med.Name
			} else {
				medNames[midStr] = "Unknown"
			}
		}
	}

	return &query.GetPatientMARResult{
		Entries: mapper.ToMAREntryResults(entries, medNames),
	}, nil
}

func (s *medicationService) GetAllMedications(ctx context.Context, q query.GetAllMedicationsQuery) (*query.GetAllMedicationsResult, error) {
	meds, err := s.medicationRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("finding all medications: %w", err)
	}
	return &query.GetAllMedicationsResult{
		Medications: mapper.ToMedicationResults(meds),
	}, nil
}
