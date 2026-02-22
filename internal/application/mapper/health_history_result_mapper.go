package mapper

import (
	"github.com/kendall/chart-attack/internal/application/common"
	"github.com/kendall/chart-attack/internal/domain/entities"
)

func ToHealthHistoryResult(h *entities.HealthHistoryEntry) common.HealthHistoryResult {
	return common.HealthHistoryResult{
		Id:           h.Id.String(),
		CreatedAt:    h.CreatedAt,
		PatientId:    h.PatientId.String(),
		Condition:    h.Condition,
		DateOccurred: h.DateOccurred,
		Description:  h.Description,
		Status:       h.Status,
	}
}

func ToHealthHistoryResults(entries []*entities.HealthHistoryEntry) []common.HealthHistoryResult {
	results := make([]common.HealthHistoryResult, 0, len(entries))
	for _, h := range entries {
		results = append(results, ToHealthHistoryResult(h))
	}
	return results
}
