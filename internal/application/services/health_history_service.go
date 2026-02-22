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

type healthHistoryService struct {
	healthHistoryRepo repositories.HealthHistoryRepository
	auditRepo         repositories.AuditRepository
	idempotencyRepo   repositories.IdempotencyRepository
}

func NewHealthHistoryService(
	healthHistoryRepo repositories.HealthHistoryRepository,
	auditRepo repositories.AuditRepository,
	idempotencyRepo repositories.IdempotencyRepository,
) appinterfaces.HealthHistoryService {
	return &healthHistoryService{
		healthHistoryRepo: healthHistoryRepo,
		auditRepo:         auditRepo,
		idempotencyRepo:   idempotencyRepo,
	}
}

func (s *healthHistoryService) AddHealthHistory(ctx context.Context, cmd command.AddHealthHistoryCommand) (*command.AddHealthHistoryResult, error) {
	if cmd.IdempotencyKey != "" {
		existing, err := s.idempotencyRepo.FindByKey(ctx, cmd.IdempotencyKey)
		if err != nil {
			return nil, fmt.Errorf("checking idempotency: %w", err)
		}
		if existing != nil {
			var result command.AddHealthHistoryResult
			if err := json.Unmarshal([]byte(existing.Response), &result); err == nil {
				return &result, nil
			}
		}
	}

	entry, err := entities.NewHealthHistoryEntry(cmd.PatientId, cmd.RecordedBy, cmd.Condition, cmd.Status)
	if err != nil {
		return nil, fmt.Errorf("creating health history entry: %w", err)
	}
	if cmd.DateOccurred != nil {
		entry.SetDateOccurred(*cmd.DateOccurred)
	}
	entry.SetDescription(cmd.Description)

	validated, err := entities.NewValidatedHealthHistoryEntry(entry)
	if err != nil {
		return nil, fmt.Errorf("validating health history entry: %w", err)
	}

	if err := s.healthHistoryRepo.Save(ctx, validated); err != nil {
		return nil, fmt.Errorf("saving health history entry: %w", err)
	}

	audit := entities.NewAuditEntry(cmd.RecordedBy, &cmd.PatientId, "CREATE", "HealthHistoryEntry", entry.Id, "{}", "", "")
	_ = s.auditRepo.Save(ctx, audit)

	result := &command.AddHealthHistoryResult{HealthHistoryId: entry.Id.String()}

	if cmd.IdempotencyKey != "" {
		responseJSON, _ := json.Marshal(result)
		record := entities.NewIdempotencyRecord(cmd.IdempotencyKey, string(responseJSON))
		_ = s.idempotencyRepo.Save(ctx, record)
	}

	return result, nil
}

func (s *healthHistoryService) GetPatientHealthHistory(ctx context.Context, q query.GetPatientHealthHistoryQuery) (*query.GetPatientHealthHistoryResult, error) {
	entries, err := s.healthHistoryRepo.FindByPatientID(ctx, q.PatientId)
	if err != nil {
		return nil, fmt.Errorf("finding health history: %w", err)
	}
	return &query.GetPatientHealthHistoryResult{
		Entries: mapper.ToHealthHistoryResults(entries),
	}, nil
}
