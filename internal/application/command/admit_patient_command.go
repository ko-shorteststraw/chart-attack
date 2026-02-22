package command

import (
	"time"

	"github.com/google/uuid"
)

type AdmitPatientCommand struct {
	IdempotencyKey string
	MRN            string
	FirstName      string
	LastName       string
	DateOfBirth    time.Time
	RoomBed        string
	Allergies      []string
	Diagnoses      []string
	CodeStatus     string
	FallRisk       bool
	IsolationType  string
	AdmittedBy     uuid.UUID
}

type AdmitPatientResult struct {
	PatientId string
}
