package query

import (
	"github.com/google/uuid"
	"github.com/kendall/chart-attack/internal/application/common"
)

type GetPatientVaccinationsQuery struct {
	PatientId uuid.UUID
}

type GetPatientVaccinationsResult struct {
	Vaccinations []common.VaccinationResult
}
