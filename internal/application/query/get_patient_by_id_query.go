package query

import (
	"github.com/google/uuid"
	"github.com/kendall/chart-attack/internal/application/common"
)

type GetPatientByIdQuery struct {
	PatientId uuid.UUID
}

type GetPatientByIdResult struct {
	Patient common.PatientResult
}
