package common

import "time"

type VitalSignResult struct {
	Id             string
	CreatedAt      time.Time
	PatientId      string
	RecordedBy     string
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
	IsAbnormal     bool
}
