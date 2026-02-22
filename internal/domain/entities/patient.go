package entities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

var validCodeStatuses = map[string]bool{
	"Full Code": true, "DNR": true, "DNI": true, "DNR/DNI": true, "Comfort Care": true,
}

var validIsolationTypes = map[string]bool{
	"None": true, "Contact": true, "Droplet": true, "Airborne": true, "Protective": true,
}

type Patient struct {
	Id              uuid.UUID
	CreatedAt       time.Time
	UpdatedAt       time.Time
	MRN             string
	FirstName       string
	LastName        string
	DateOfBirth     time.Time
	RoomBed         string
	Allergies       []string
	Diagnoses       []string
	CodeStatus      string
	FallRisk        bool
	IsolationType   string
	AdmitDate       time.Time
	DischargeDate   *time.Time
	AssignedNurseId *uuid.UUID
}

func NewPatient(mrn, firstName, lastName string, dob time.Time, roomBed string, allergies, diagnoses []string, codeStatus, isolationType string, fallRisk bool) (*Patient, error) {
	now := time.Now().UTC()
	if allergies == nil {
		allergies = []string{}
	}
	if diagnoses == nil {
		diagnoses = []string{}
	}
	p := &Patient{
		Id:            uuid.Must(uuid.NewV7()),
		CreatedAt:     now,
		UpdatedAt:     now,
		MRN:           mrn,
		FirstName:     firstName,
		LastName:      lastName,
		DateOfBirth:   dob,
		RoomBed:       roomBed,
		Allergies:     allergies,
		Diagnoses:     diagnoses,
		CodeStatus:    codeStatus,
		FallRisk:      fallRisk,
		IsolationType: isolationType,
		AdmitDate:     now,
	}
	if err := p.validate(); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Patient) validate() error {
	if p.MRN == "" {
		return fmt.Errorf("MRN is required")
	}
	if p.FirstName == "" {
		return fmt.Errorf("first name is required")
	}
	if p.LastName == "" {
		return fmt.Errorf("last name is required")
	}
	if p.RoomBed == "" {
		return fmt.Errorf("room/bed is required")
	}
	if !validCodeStatuses[p.CodeStatus] {
		return fmt.Errorf("invalid code status: %s", p.CodeStatus)
	}
	if !validIsolationTypes[p.IsolationType] {
		return fmt.Errorf("invalid isolation type: %s", p.IsolationType)
	}
	return nil
}

func (p *Patient) UpdateCodeStatus(status string) error {
	p.CodeStatus = status
	p.UpdatedAt = time.Now().UTC()
	return p.validate()
}

func (p *Patient) UpdateRoomBed(roomBed string) error {
	p.RoomBed = roomBed
	p.UpdatedAt = time.Now().UTC()
	return p.validate()
}

func (p *Patient) UpdateFallRisk(fallRisk bool) {
	p.FallRisk = fallRisk
	p.UpdatedAt = time.Now().UTC()
}

func (p *Patient) UpdateIsolationType(isoType string) error {
	p.IsolationType = isoType
	p.UpdatedAt = time.Now().UTC()
	return p.validate()
}

func (p *Patient) AssignNurse(nurseId uuid.UUID) {
	p.AssignedNurseId = &nurseId
	p.UpdatedAt = time.Now().UTC()
}

func (p *Patient) Discharge() {
	now := time.Now().UTC()
	p.DischargeDate = &now
	p.UpdatedAt = now
}

func (p *Patient) IsAdmitted() bool {
	return p.DischargeDate == nil
}

func (p *Patient) FullName() string {
	return p.LastName + ", " + p.FirstName
}

func (p *Patient) AddDiagnosis(d string) {
	for _, existing := range p.Diagnoses {
		if existing == d {
			return
		}
	}
	p.Diagnoses = append(p.Diagnoses, d)
	p.UpdatedAt = time.Now().UTC()
}

func (p *Patient) RemoveDiagnosis(d string) {
	for i, existing := range p.Diagnoses {
		if existing == d {
			p.Diagnoses = append(p.Diagnoses[:i], p.Diagnoses[i+1:]...)
			p.UpdatedAt = time.Now().UTC()
			return
		}
	}
}
