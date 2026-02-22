package common

import "time"

type TaskResult struct {
	Id          string
	PatientId   string
	AssignedTo  string
	Title       string
	Category    string
	DueAt       time.Time
	CompletedAt *time.Time
	Priority    string
	Recurring   bool
	Notes       string
	IsCompleted bool
}
