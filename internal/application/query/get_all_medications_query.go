package query

import "github.com/kendall/chart-attack/internal/application/common"

type GetAllMedicationsQuery struct{}

type GetAllMedicationsResult struct {
	Medications []common.MedicationResult
}
