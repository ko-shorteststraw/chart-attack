package query

import (
	"github.com/google/uuid"
	"github.com/kendall/chart-attack/internal/application/common"
)

type GetPatientMARQuery struct {
	PatientId uuid.UUID
}

type GetPatientMARResult struct {
	Entries []common.MAREntryResult
}
