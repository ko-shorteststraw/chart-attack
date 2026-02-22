package command

import (
	"time"

	"github.com/google/uuid"
)

type AdministerMedicationCommand struct {
	IdempotencyKey string
	PatientId      uuid.UUID
	MedicationId   uuid.UUID
	ScheduledTime  time.Time
	Dose           string
	Route          string
	AdministeredBy uuid.UUID
}

type AdministerMedicationResult struct {
	MAREntryId string
}
