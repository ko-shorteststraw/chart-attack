package query

import (
	"github.com/google/uuid"
	"github.com/kendall/chart-attack/internal/application/common"
)

type GetPatientHealthHistoryQuery struct {
	PatientId uuid.UUID
}

type GetPatientHealthHistoryResult struct {
	Entries []common.HealthHistoryResult
}
