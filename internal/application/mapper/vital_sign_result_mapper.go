package mapper

import (
	"github.com/kendall/chart-attack/internal/application/common"
	"github.com/kendall/chart-attack/internal/domain/entities"
)

func ToVitalSignResult(vs *entities.VitalSign) common.VitalSignResult {
	return common.VitalSignResult{
		Id:             vs.Id.String(),
		CreatedAt:      vs.CreatedAt,
		PatientId:      vs.PatientId.String(),
		RecordedBy:     vs.RecordedBy.String(),
		SystolicBP:     vs.SystolicBP,
		DiastolicBP:    vs.DiastolicBP,
		HeartRate:      vs.HeartRate,
		Temperature:    vs.Temperature,
		TempRoute:      vs.TempRoute,
		OxygenSat:      vs.OxygenSat,
		Respirations:   vs.Respirations,
		PainLevel:      vs.PainLevel,
		SupplementalO2: vs.SupplementalO2,
		O2FlowRate:     vs.O2FlowRate,
		Position:       vs.Position,
		Notes:          vs.Notes,
		IsAbnormal:     vs.IsAbnormal(),
	}
}

func ToVitalSignResults(vitals []*entities.VitalSign) []common.VitalSignResult {
	results := make([]common.VitalSignResult, 0, len(vitals))
	for _, vs := range vitals {
		results = append(results, ToVitalSignResult(vs))
	}
	return results
}
