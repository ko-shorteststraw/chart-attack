package entities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

var validTempRoutes = map[string]bool{
	"": true, "Oral": true, "Tympanic": true, "Axillary": true, "Rectal": true, "Temporal": true,
}

var validPositions = map[string]bool{
	"": true, "Sitting": true, "Standing": true, "Supine": true, "Left lateral": true,
}

type VitalSign struct {
	Id             uuid.UUID
	CreatedAt      time.Time
	PatientId      uuid.UUID
	RecordedBy     uuid.UUID
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
}

func NewVitalSign(patientId, recordedBy uuid.UUID) (*VitalSign, error) {
	now := time.Now().UTC()
	vs := &VitalSign{
		Id:        uuid.Must(uuid.NewV7()),
		CreatedAt: now,
		PatientId: patientId,
		RecordedBy: recordedBy,
	}
	return vs, nil
}

func (vs *VitalSign) validate() error {
	if vs.PatientId == uuid.Nil {
		return fmt.Errorf("patient ID is required")
	}
	if vs.RecordedBy == uuid.Nil {
		return fmt.Errorf("recorded by user ID is required")
	}
	if vs.SystolicBP != nil && (*vs.SystolicBP < 0 || *vs.SystolicBP > 400) {
		return fmt.Errorf("systolic BP out of range: %d", *vs.SystolicBP)
	}
	if vs.DiastolicBP != nil && (*vs.DiastolicBP < 0 || *vs.DiastolicBP > 300) {
		return fmt.Errorf("diastolic BP out of range: %d", *vs.DiastolicBP)
	}
	if vs.HeartRate != nil && (*vs.HeartRate < 0 || *vs.HeartRate > 400) {
		return fmt.Errorf("heart rate out of range: %d", *vs.HeartRate)
	}
	if vs.Temperature != nil && (*vs.Temperature < 80.0 || *vs.Temperature > 115.0) {
		return fmt.Errorf("temperature out of range: %.1f", *vs.Temperature)
	}
	if !validTempRoutes[vs.TempRoute] {
		return fmt.Errorf("invalid temp route: %s", vs.TempRoute)
	}
	if vs.OxygenSat != nil && (*vs.OxygenSat < 0 || *vs.OxygenSat > 100) {
		return fmt.Errorf("oxygen sat out of range: %d", *vs.OxygenSat)
	}
	if vs.Respirations != nil && (*vs.Respirations < 0 || *vs.Respirations > 100) {
		return fmt.Errorf("respirations out of range: %d", *vs.Respirations)
	}
	if vs.PainLevel != nil && (*vs.PainLevel < 0 || *vs.PainLevel > 10) {
		return fmt.Errorf("pain level out of range: %d", *vs.PainLevel)
	}
	if !validPositions[vs.Position] {
		return fmt.Errorf("invalid position: %s", vs.Position)
	}
	return nil
}

func (vs *VitalSign) SetBloodPressure(systolic, diastolic int) error {
	vs.SystolicBP = &systolic
	vs.DiastolicBP = &diastolic
	return vs.validate()
}

func (vs *VitalSign) SetHeartRate(hr int) error {
	vs.HeartRate = &hr
	return vs.validate()
}

func (vs *VitalSign) SetTemperature(temp float64, route string) error {
	vs.Temperature = &temp
	vs.TempRoute = route
	return vs.validate()
}

func (vs *VitalSign) SetOxygenSat(sat int) error {
	vs.OxygenSat = &sat
	return vs.validate()
}

func (vs *VitalSign) SetRespirations(resp int) error {
	vs.Respirations = &resp
	return vs.validate()
}

func (vs *VitalSign) SetPainLevel(level int) error {
	vs.PainLevel = &level
	return vs.validate()
}

func (vs *VitalSign) SetSupplementalO2(on bool, flowRate *float64) {
	vs.SupplementalO2 = on
	vs.O2FlowRate = flowRate
}

func (vs *VitalSign) SetPosition(position string) error {
	vs.Position = position
	return vs.validate()
}

func (vs *VitalSign) SetNotes(notes string) {
	vs.Notes = notes
}

// IsAbnormal checks if any vital sign values are outside normal ranges.
func (vs *VitalSign) IsAbnormal() bool {
	if vs.SystolicBP != nil && (*vs.SystolicBP > 180 || *vs.SystolicBP < 90) {
		return true
	}
	if vs.DiastolicBP != nil && (*vs.DiastolicBP > 120 || *vs.DiastolicBP < 60) {
		return true
	}
	if vs.HeartRate != nil && (*vs.HeartRate > 120 || *vs.HeartRate < 50) {
		return true
	}
	if vs.Temperature != nil && (*vs.Temperature > 100.4 || *vs.Temperature < 96.0) {
		return true
	}
	if vs.OxygenSat != nil && *vs.OxygenSat < 92 {
		return true
	}
	if vs.Respirations != nil && (*vs.Respirations > 24 || *vs.Respirations < 10) {
		return true
	}
	return false
}
