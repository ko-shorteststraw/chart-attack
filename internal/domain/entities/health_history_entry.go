package entities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

var validHistoryStatuses = map[string]bool{
	"Active": true, "Resolved": true, "Chronic": true,
}

type HealthHistoryEntry struct {
	Id           uuid.UUID
	CreatedAt    time.Time
	PatientId    uuid.UUID
	RecordedBy   uuid.UUID
	Condition    string
	DateOccurred *time.Time
	Description  string
	Status       string
}

func NewHealthHistoryEntry(patientId, recordedBy uuid.UUID, condition, status string) (*HealthHistoryEntry, error) {
	now := time.Now().UTC()
	h := &HealthHistoryEntry{
		Id:         uuid.Must(uuid.NewV7()),
		CreatedAt:  now,
		PatientId:  patientId,
		RecordedBy: recordedBy,
		Condition:  condition,
		Status:     status,
	}
	if err := h.validate(); err != nil {
		return nil, err
	}
	return h, nil
}

func (h *HealthHistoryEntry) validate() error {
	if h.PatientId == uuid.Nil {
		return fmt.Errorf("patient ID is required")
	}
	if h.RecordedBy == uuid.Nil {
		return fmt.Errorf("recorded by user ID is required")
	}
	if h.Condition == "" {
		return fmt.Errorf("condition is required")
	}
	if !validHistoryStatuses[h.Status] {
		return fmt.Errorf("invalid status: %s", h.Status)
	}
	return nil
}

func (h *HealthHistoryEntry) SetDateOccurred(t time.Time) {
	h.DateOccurred = &t
}

func (h *HealthHistoryEntry) SetDescription(desc string) {
	h.Description = desc
}
