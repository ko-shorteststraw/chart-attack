package entities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type VaccinationRecord struct {
	Id               uuid.UUID
	CreatedAt        time.Time
	PatientId        uuid.UUID
	RecordedBy       uuid.UUID
	VaccineName      string
	DateAdministered time.Time
	LotNumber        string
	Site             string
	Notes            string
}

func NewVaccinationRecord(patientId, recordedBy uuid.UUID, vaccineName string, dateAdministered time.Time) (*VaccinationRecord, error) {
	now := time.Now().UTC()
	v := &VaccinationRecord{
		Id:               uuid.Must(uuid.NewV7()),
		CreatedAt:        now,
		PatientId:        patientId,
		RecordedBy:       recordedBy,
		VaccineName:      vaccineName,
		DateAdministered: dateAdministered,
	}
	if err := v.validate(); err != nil {
		return nil, err
	}
	return v, nil
}

func (v *VaccinationRecord) validate() error {
	if v.PatientId == uuid.Nil {
		return fmt.Errorf("patient ID is required")
	}
	if v.RecordedBy == uuid.Nil {
		return fmt.Errorf("recorded by user ID is required")
	}
	if v.VaccineName == "" {
		return fmt.Errorf("vaccine name is required")
	}
	return nil
}

func (v *VaccinationRecord) SetLotNumber(lot string) {
	v.LotNumber = lot
}

func (v *VaccinationRecord) SetSite(site string) {
	v.Site = site
}

func (v *VaccinationRecord) SetNotes(notes string) {
	v.Notes = notes
}
