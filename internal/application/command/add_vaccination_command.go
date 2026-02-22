package command

import (
	"time"

	"github.com/google/uuid"
)

type AddVaccinationCommand struct {
	IdempotencyKey   string
	PatientId        uuid.UUID
	RecordedBy       uuid.UUID
	VaccineName      string
	DateAdministered time.Time
	LotNumber        string
	Site             string
	Notes            string
}

type AddVaccinationResult struct {
	VaccinationId string
}
