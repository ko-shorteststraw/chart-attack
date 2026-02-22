package command

import "github.com/google/uuid"

type AddEmergencyContactCommand struct {
	IdempotencyKey string
	PatientId      uuid.UUID
	Name           string
	Relationship   string
	Phone          string
	Email          string
	IsPrimary      bool
	AddedBy        uuid.UUID
}

type AddEmergencyContactResult struct {
	ContactId string
}
