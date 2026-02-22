package entities

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPatient(t *testing.T) {
	dob := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)
	p, err := NewPatient("MRN001", "John", "Doe", dob, "401-A",
		[]string{"Penicillin"}, []string{"CHF"}, "Full Code", "None", false)

	require.NoError(t, err)
	assert.Equal(t, "MRN001", p.MRN)
	assert.Equal(t, "John", p.FirstName)
	assert.Equal(t, "Doe", p.LastName)
	assert.Equal(t, "401-A", p.RoomBed)
	assert.Equal(t, []string{"Penicillin"}, p.Allergies)
	assert.Equal(t, "Full Code", p.CodeStatus)
	assert.False(t, p.FallRisk)
	assert.True(t, p.IsAdmitted())
	assert.Equal(t, "Doe, John", p.FullName())
	assert.NotEmpty(t, p.Id)
}

func TestNewPatient_ValidationErrors(t *testing.T) {
	dob := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)

	_, err := NewPatient("", "John", "Doe", dob, "401-A", nil, nil, "Full Code", "None", false)
	assert.Error(t, err)

	_, err = NewPatient("MRN001", "", "Doe", dob, "401-A", nil, nil, "Full Code", "None", false)
	assert.Error(t, err)

	_, err = NewPatient("MRN001", "John", "Doe", dob, "401-A", nil, nil, "Invalid", "None", false)
	assert.Error(t, err)

	_, err = NewPatient("MRN001", "John", "Doe", dob, "401-A", nil, nil, "Full Code", "Invalid", false)
	assert.Error(t, err)
}

func TestPatient_Discharge(t *testing.T) {
	dob := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)
	p, _ := NewPatient("MRN001", "John", "Doe", dob, "401-A", nil, nil, "Full Code", "None", false)

	assert.True(t, p.IsAdmitted())
	p.Discharge()
	assert.False(t, p.IsAdmitted())
	assert.NotNil(t, p.DischargeDate)
}

func TestPatient_UpdateCodeStatus(t *testing.T) {
	dob := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)
	p, _ := NewPatient("MRN001", "John", "Doe", dob, "401-A", nil, nil, "Full Code", "None", false)

	err := p.UpdateCodeStatus("DNR")
	assert.NoError(t, err)
	assert.Equal(t, "DNR", p.CodeStatus)

	err = p.UpdateCodeStatus("Invalid")
	assert.Error(t, err)
}

func TestValidatedPatient(t *testing.T) {
	dob := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)
	p, _ := NewPatient("MRN001", "John", "Doe", dob, "401-A", nil, nil, "Full Code", "None", false)

	vp, err := NewValidatedPatient(p)
	require.NoError(t, err)
	assert.Equal(t, p, vp.Patient())
}
