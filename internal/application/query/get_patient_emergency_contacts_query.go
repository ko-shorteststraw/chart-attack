package query

import (
	"github.com/google/uuid"
	"github.com/kendall/chart-attack/internal/application/common"
)

type GetPatientEmergencyContactsQuery struct {
	PatientId uuid.UUID
}

type GetPatientEmergencyContactsResult struct {
	Contacts []common.EmergencyContactResult
}
