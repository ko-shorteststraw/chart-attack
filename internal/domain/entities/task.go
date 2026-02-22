package entities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

var validTaskCategories = map[string]bool{
	"Medication": true, "Vital": true, "Turn": true, "Assessment": true, "Lab": true, "Custom": true,
}

var validTaskPriorities = map[string]bool{
	"Routine": true, "Urgent": true, "STAT": true,
}

type Task struct {
	Id            uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	PatientId     uuid.UUID
	AssignedTo    uuid.UUID
	Title         string
	Category      string
	DueAt         time.Time
	CompletedAt   *time.Time
	CompletedBy   *uuid.UUID
	Priority      string
	Recurring     bool
	RecurInterval *string
	Notes         string
}

func NewTask(patientId, assignedTo uuid.UUID, title, category string, dueAt time.Time, priority string) (*Task, error) {
	now := time.Now().UTC()
	t := &Task{
		Id:         uuid.Must(uuid.NewV7()),
		CreatedAt:  now,
		UpdatedAt:  now,
		PatientId:  patientId,
		AssignedTo: assignedTo,
		Title:      title,
		Category:   category,
		DueAt:      dueAt,
		Priority:   priority,
	}
	if err := t.validate(); err != nil {
		return nil, err
	}
	return t, nil
}

func (t *Task) validate() error {
	if t.PatientId == uuid.Nil {
		return fmt.Errorf("patient ID is required")
	}
	if t.AssignedTo == uuid.Nil {
		return fmt.Errorf("assigned to user ID is required")
	}
	if t.Title == "" {
		return fmt.Errorf("task title is required")
	}
	if !validTaskCategories[t.Category] {
		return fmt.Errorf("invalid task category: %s", t.Category)
	}
	if !validTaskPriorities[t.Priority] {
		return fmt.Errorf("invalid task priority: %s", t.Priority)
	}
	return nil
}

func (t *Task) Complete(completedBy uuid.UUID) {
	now := time.Now().UTC()
	t.CompletedAt = &now
	t.CompletedBy = &completedBy
	t.UpdatedAt = now
}

func (t *Task) IsCompleted() bool {
	return t.CompletedAt != nil
}

func (t *Task) SetNotes(notes string) {
	t.Notes = notes
	t.UpdatedAt = time.Now().UTC()
}
