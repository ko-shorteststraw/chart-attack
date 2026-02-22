package mapper

import (
	"github.com/kendall/chart-attack/internal/application/common"
	"github.com/kendall/chart-attack/internal/domain/entities"
)

func ToMedicationResult(m *entities.Medication) common.MedicationResult {
	return common.MedicationResult{
		Id:           m.Id.String(),
		Name:         m.Name,
		BrandName:    m.BrandName,
		DrugClass:    m.DrugClass,
		DefaultDose:  m.DefaultDose,
		DefaultRoute: m.DefaultRoute,
		Frequency:    m.Frequency,
		HighAlert:    m.HighAlert,
	}
}

func ToMedicationResults(meds []*entities.Medication) []common.MedicationResult {
	results := make([]common.MedicationResult, 0, len(meds))
	for _, m := range meds {
		results = append(results, ToMedicationResult(m))
	}
	return results
}
