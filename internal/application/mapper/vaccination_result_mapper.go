package mapper

import (
	"github.com/kendall/chart-attack/internal/application/common"
	"github.com/kendall/chart-attack/internal/domain/entities"
)

func ToVaccinationResult(v *entities.VaccinationRecord) common.VaccinationResult {
	return common.VaccinationResult{
		Id:               v.Id.String(),
		CreatedAt:        v.CreatedAt,
		PatientId:        v.PatientId.String(),
		VaccineName:      v.VaccineName,
		DateAdministered: v.DateAdministered,
		LotNumber:        v.LotNumber,
		Site:             v.Site,
		Notes:            v.Notes,
	}
}

func ToVaccinationResults(records []*entities.VaccinationRecord) []common.VaccinationResult {
	results := make([]common.VaccinationResult, 0, len(records))
	for _, v := range records {
		results = append(results, ToVaccinationResult(v))
	}
	return results
}
