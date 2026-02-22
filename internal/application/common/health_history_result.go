package common

import "time"

type HealthHistoryResult struct {
	Id           string
	CreatedAt    time.Time
	PatientId    string
	Condition    string
	DateOccurred *time.Time
	Description  string
	Status       string
}
