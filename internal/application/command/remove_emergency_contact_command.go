package command

import "github.com/google/uuid"

type RemoveEmergencyContactCommand struct {
	ContactId uuid.UUID
	RemovedBy uuid.UUID
}
