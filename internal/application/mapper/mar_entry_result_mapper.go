package mapper

import (
	"github.com/kendall/chart-attack/internal/application/common"
	"github.com/kendall/chart-attack/internal/domain/entities"
)

func ToMAREntryResult(m *entities.MAREntry, medName string) common.MAREntryResult {
	return common.MAREntryResult{
		Id:             m.Id.String(),
		CreatedAt:      m.CreatedAt,
		PatientId:      m.PatientId.String(),
		MedicationId:   m.MedicationId.String(),
		MedicationName: medName,
		ScheduledTime:  m.ScheduledTime,
		AdministeredAt: m.AdministeredAt,
		Status:         m.Status,
		Dose:           m.Dose,
		Route:          m.Route,
		Site:           m.Site,
		HoldReason:     m.HoldReason,
		Notes:          m.Notes,
	}
}

func ToMAREntryResults(entries []*entities.MAREntry, medNames map[string]string) []common.MAREntryResult {
	results := make([]common.MAREntryResult, 0, len(entries))
	for _, e := range entries {
		name := medNames[e.MedicationId.String()]
		results = append(results, ToMAREntryResult(e, name))
	}
	return results
}
