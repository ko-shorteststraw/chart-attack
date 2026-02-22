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

type vaccinationService struct {
	vaccinationRepo repositories.VaccinationRepository
	auditRepo       repositories.AuditRepository
	idempotencyRepo repositories.IdempotencyRepository
}

func NewVaccinationService(
	vaccinationRepo repositories.VaccinationRepository,
	auditRepo repositories.AuditRepository,
	idempotencyRepo repositories.IdempotencyRepository,
) appinterfaces.VaccinationService {
	return &vaccinationService{
		vaccinationRepo: vaccinationRepo,
		auditRepo:       auditRepo,
		idempotencyRepo: idempotencyRepo,
	}
}

func (s *vaccinationService) AddVaccination(ctx context.Context, cmd command.AddVaccinationCommand) (*command.AddVaccinationResult, error) {
	if cmd.IdempotencyKey != "" {
		existing, err := s.idempotencyRepo.FindByKey(ctx, cmd.IdempotencyKey)
		if err != nil {
			return nil, fmt.Errorf("checking idempotency: %w", err)
		}
		if existing != nil {
			var result command.AddVaccinationResult
			if err := json.Unmarshal([]byte(existing.Response), &result); err == nil {
				return &result, nil
			}
		}
	}

	vac, err := entities.NewVaccinationRecord(cmd.PatientId, cmd.RecordedBy, cmd.VaccineName, cmd.DateAdministered)
	if err != nil {
		return nil, fmt.Errorf("creating vaccination record: %w", err)
	}
	vac.SetLotNumber(cmd.LotNumber)
	vac.SetSite(cmd.Site)
	vac.SetNotes(cmd.Notes)

	validated, err := entities.NewValidatedVaccinationRecord(vac)
	if err != nil {
		return nil, fmt.Errorf("validating vaccination record: %w", err)
	}

	if err := s.vaccinationRepo.Save(ctx, validated); err != nil {
		return nil, fmt.Errorf("saving vaccination record: %w", err)
	}

	audit := entities.NewAuditEntry(cmd.RecordedBy, &cmd.PatientId, "CREATE", "VaccinationRecord", vac.Id, "{}", "", "")
	_ = s.auditRepo.Save(ctx, audit)

	result := &command.AddVaccinationResult{VaccinationId: vac.Id.String()}

	if cmd.IdempotencyKey != "" {
		responseJSON, _ := json.Marshal(result)
		record := entities.NewIdempotencyRecord(cmd.IdempotencyKey, string(responseJSON))
		_ = s.idempotencyRepo.Save(ctx, record)
	}

	return result, nil
}

func (s *vaccinationService) GetPatientVaccinations(ctx context.Context, q query.GetPatientVaccinationsQuery) (*query.GetPatientVaccinationsResult, error) {
	records, err := s.vaccinationRepo.FindByPatientID(ctx, q.PatientId)
	if err != nil {
		return nil, fmt.Errorf("finding vaccinations: %w", err)
	}
	return &query.GetPatientVaccinationsResult{
		Vaccinations: mapper.ToVaccinationResults(records),
	}, nil
}
