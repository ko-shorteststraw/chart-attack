package command

import (
	"time"

	"github.com/google/uuid"
)

type AddHealthHistoryCommand struct {
	IdempotencyKey string
	PatientId      uuid.UUID
	RecordedBy     uuid.UUID
	Condition      string
	DateOccurred   *time.Time
	Description    string
	Status         string
}

type AddHealthHistoryResult struct {
	HealthHistoryId string
}
