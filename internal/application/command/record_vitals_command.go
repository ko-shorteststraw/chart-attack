package command

import "github.com/google/uuid"

type RecordVitalsCommand struct {
	IdempotencyKey string
	PatientId      uuid.UUID
	RecordedBy     uuid.UUID
	SystolicBP     *int
	DiastolicBP    *int
	HeartRate      *int
	Temperature    *float64
	TempRoute      string
	OxygenSat      *int
	Respirations   *int
	PainLevel      *int
	SupplementalO2 bool
	O2FlowRate     *float64
	Position       string
	Notes          string
}

type RecordVitalsResult struct {
	VitalSignId string
	IsAbnormal  bool
}
