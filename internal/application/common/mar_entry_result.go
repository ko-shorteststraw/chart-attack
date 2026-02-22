package common

import "time"

type MAREntryResult struct {
	Id             string
	CreatedAt      time.Time
	PatientId      string
	MedicationId   string
	MedicationName string
	ScheduledTime  time.Time
	AdministeredAt *time.Time
	Status         string
	Dose           string
	Route          string
	Site           string
	HoldReason     string
	Notes          string
}
