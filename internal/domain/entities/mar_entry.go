package entities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

var validMARStatuses = map[string]bool{
	"Scheduled": true, "Given": true, "Held": true, "Refused": true, "Omitted": true,
}

var validRoutes = map[string]bool{
	"PO": true, "IV": true, "IM": true, "SQ": true, "SL": true, "PR": true, "Topical": true, "Inhaled": true, "": true,
}

type MAREntry struct {
	Id             uuid.UUID
	CreatedAt      time.Time
	PatientId      uuid.UUID
	MedicationId   uuid.UUID
	ScheduledTime  time.Time
	AdministeredAt *time.Time
	AdministeredBy *uuid.UUID
	Status         string
	Dose           string
	Route          string
	Site           string
	HoldReason     string
	Notes          string
}

func NewMAREntry(patientId, medicationId uuid.UUID, scheduledTime time.Time, dose, route string) (*MAREntry, error) {
	now := time.Now().UTC()
	m := &MAREntry{
		Id:            uuid.Must(uuid.NewV7()),
		CreatedAt:     now,
		PatientId:     patientId,
		MedicationId:  medicationId,
		ScheduledTime: scheduledTime,
		Status:        "Scheduled",
		Dose:          dose,
		Route:         route,
	}
	if err := m.validate(); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *MAREntry) validate() error {
	if m.PatientId == uuid.Nil {
		return fmt.Errorf("patient ID is required")
	}
	if m.MedicationId == uuid.Nil {
		return fmt.Errorf("medication ID is required")
	}
	if !validMARStatuses[m.Status] {
		return fmt.Errorf("invalid MAR status: %s", m.Status)
	}
	if !validRoutes[m.Route] {
		return fmt.Errorf("invalid route: %s", m.Route)
	}
	return nil
}

func (m *MAREntry) Administer(administeredBy uuid.UUID) error {
	now := time.Now().UTC()
	m.AdministeredAt = &now
	m.AdministeredBy = &administeredBy
	m.Status = "Given"
	return m.validate()
}

func (m *MAREntry) Hold(reason string) error {
	m.Status = "Held"
	m.HoldReason = reason
	return m.validate()
}

func (m *MAREntry) SetNotes(notes string) {
	m.Notes = notes
}
