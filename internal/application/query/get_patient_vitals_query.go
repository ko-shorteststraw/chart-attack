package query

import (
	"github.com/google/uuid"
	"github.com/kendall/chart-attack/internal/application/common"
)

type GetPatientVitalsQuery struct {
	PatientId uuid.UUID
}

type GetPatientVitalsResult struct {
	Vitals []common.VitalSignResult
}
