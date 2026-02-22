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

type chartingService struct {
	vitalSignRepo   repositories.VitalSignRepository
	auditRepo       repositories.AuditRepository
	idempotencyRepo repositories.IdempotencyRepository
}

func NewChartingService(
	vitalSignRepo repositories.VitalSignRepository,
	auditRepo repositories.AuditRepository,
	idempotencyRepo repositories.IdempotencyRepository,
) appinterfaces.ChartingService {
	return &chartingService{
		vitalSignRepo:   vitalSignRepo,
		auditRepo:       auditRepo,
		idempotencyRepo: idempotencyRepo,
	}
}

func (s *chartingService) RecordVitals(ctx context.Context, cmd command.RecordVitalsCommand) (*command.RecordVitalsResult, error) {
	// Check idempotency
	if cmd.IdempotencyKey != "" {
		existing, err := s.idempotencyRepo.FindByKey(ctx, cmd.IdempotencyKey)
		if err != nil {
			return nil, fmt.Errorf("checking idempotency: %w", err)
		}
		if existing != nil {
			var result command.RecordVitalsResult
			if err := json.Unmarshal([]byte(existing.Response), &result); err == nil {
				return &result, nil
			}
		}
	}

	vs, err := entities.NewVitalSign(cmd.PatientId, cmd.RecordedBy)
	if err != nil {
		return nil, fmt.Errorf("creating vital sign: %w", err)
	}

	// Set optional fields
	if cmd.SystolicBP != nil && cmd.DiastolicBP != nil {
		if err := vs.SetBloodPressure(*cmd.SystolicBP, *cmd.DiastolicBP); err != nil {
			return nil, fmt.Errorf("setting blood pressure: %w", err)
		}
	}
	if cmd.HeartRate != nil {
		if err := vs.SetHeartRate(*cmd.HeartRate); err != nil {
			return nil, fmt.Errorf("setting heart rate: %w", err)
		}
	}
	if cmd.Temperature != nil {
		if err := vs.SetTemperature(*cmd.Temperature, cmd.TempRoute); err != nil {
			return nil, fmt.Errorf("setting temperature: %w", err)
		}
	}
	if cmd.OxygenSat != nil {
		if err := vs.SetOxygenSat(*cmd.OxygenSat); err != nil {
			return nil, fmt.Errorf("setting oxygen sat: %w", err)
		}
	}
	if cmd.Respirations != nil {
		if err := vs.SetRespirations(*cmd.Respirations); err != nil {
			return nil, fmt.Errorf("setting respirations: %w", err)
		}
	}
	if cmd.PainLevel != nil {
		if err := vs.SetPainLevel(*cmd.PainLevel); err != nil {
			return nil, fmt.Errorf("setting pain level: %w", err)
		}
	}
	vs.SetSupplementalO2(cmd.SupplementalO2, cmd.O2FlowRate)
	if cmd.Position != "" {
		if err := vs.SetPosition(cmd.Position); err != nil {
			return nil, fmt.Errorf("setting position: %w", err)
		}
	}
	vs.SetNotes(cmd.Notes)

	validated, err := entities.NewValidatedVitalSign(vs)
	if err != nil {
		return nil, fmt.Errorf("validating vital sign: %w", err)
	}

	if err := s.vitalSignRepo.Save(ctx, validated); err != nil {
		return nil, fmt.Errorf("saving vital sign: %w", err)
	}

	// Audit trail
	audit := entities.NewAuditEntry(cmd.RecordedBy, &cmd.PatientId, "CREATE", "VitalSign", vs.Id, "{}", "", "")
	if err := s.auditRepo.Save(ctx, audit); err != nil {
		return nil, fmt.Errorf("saving audit entry: %w", err)
	}

	result := &command.RecordVitalsResult{
		VitalSignId: vs.Id.String(),
		IsAbnormal:  vs.IsAbnormal(),
	}

	// Save idempotency record
	if cmd.IdempotencyKey != "" {
		responseJSON, _ := json.Marshal(result)
		record := entities.NewIdempotencyRecord(cmd.IdempotencyKey, string(responseJSON))
		_ = s.idempotencyRepo.Save(ctx, record)
	}

	return result, nil
}

func (s *chartingService) GetPatientVitals(ctx context.Context, q query.GetPatientVitalsQuery) (*query.GetPatientVitalsResult, error) {
	vitals, err := s.vitalSignRepo.FindByPatientID(ctx, q.PatientId)
	if err != nil {
		return nil, fmt.Errorf("finding patient vitals: %w", err)
	}
	return &query.GetPatientVitalsResult{
		Vitals: mapper.ToVitalSignResults(vitals),
	}, nil
}
