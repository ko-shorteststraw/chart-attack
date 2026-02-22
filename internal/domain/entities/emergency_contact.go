package entities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type EmergencyContact struct {
	Id           uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
	PatientId    uuid.UUID
	Name         string
	Relationship string
	Phone        string
	Email        string
	IsPrimary    bool
}

func NewEmergencyContact(patientId uuid.UUID, name, relationship, phone, email string, isPrimary bool) (*EmergencyContact, error) {
	now := time.Now().UTC()
	ec := &EmergencyContact{
		Id:           uuid.Must(uuid.NewV7()),
		CreatedAt:    now,
		UpdatedAt:    now,
		PatientId:    patientId,
		Name:         name,
		Relationship: relationship,
		Phone:        phone,
		Email:        email,
		IsPrimary:    isPrimary,
	}
	if err := ec.validate(); err != nil {
		return nil, err
	}
	return ec, nil
}

func (ec *EmergencyContact) validate() error {
	if ec.PatientId == uuid.Nil {
		return fmt.Errorf("patient ID is required")
	}
	if ec.Name == "" {
		return fmt.Errorf("contact name is required")
	}
	return nil
}

func (ec *EmergencyContact) UpdatePhone(phone string) {
	ec.Phone = phone
	ec.UpdatedAt = time.Now().UTC()
}

func (ec *EmergencyContact) UpdateEmail(email string) {
	ec.Email = email
	ec.UpdatedAt = time.Now().UTC()
}
