package command

import "github.com/google/uuid"

type UpdateDiagnosesCommand struct {
	PatientId uuid.UUID
	Action    string // "add" or "remove"
	Diagnosis string
	UpdatedBy uuid.UUID
}
