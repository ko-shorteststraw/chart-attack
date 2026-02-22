package command

import "github.com/google/uuid"

type CompleteTaskCommand struct {
	TaskId      uuid.UUID
	CompletedBy uuid.UUID
}
