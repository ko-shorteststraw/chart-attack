package command

import (
	"time"

	"github.com/google/uuid"
)

type CreateTaskCommand struct {
	IdempotencyKey string
	PatientId      uuid.UUID
	AssignedTo     uuid.UUID
	Title          string
	Category       string
	DueAt          time.Time
	Priority       string
	Notes          string
	CreatedBy      uuid.UUID
}

type CreateTaskResult struct {
	TaskId string
}
