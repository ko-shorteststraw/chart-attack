package entities

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewVitalSign(t *testing.T) {
	patientId := uuid.Must(uuid.NewV7())
	nurseId := uuid.Must(uuid.NewV7())

	vs, err := NewVitalSign(patientId, nurseId)
	require.NoError(t, err)
	assert.Equal(t, patientId, vs.PatientId)
	assert.Equal(t, nurseId, vs.RecordedBy)
	assert.NotEmpty(t, vs.Id)
}

func TestVitalSign_SetBloodPressure(t *testing.T) {
	vs, _ := NewVitalSign(uuid.Must(uuid.NewV7()), uuid.Must(uuid.NewV7()))

	err := vs.SetBloodPressure(120, 80)
	assert.NoError(t, err)
	assert.Equal(t, 120, *vs.SystolicBP)
	assert.Equal(t, 80, *vs.DiastolicBP)

	err = vs.SetBloodPressure(500, 80)
	assert.Error(t, err)
}

func TestVitalSign_SetTemperature(t *testing.T) {
	vs, _ := NewVitalSign(uuid.Must(uuid.NewV7()), uuid.Must(uuid.NewV7()))

	err := vs.SetTemperature(98.6, "Oral")
	assert.NoError(t, err)
	assert.Equal(t, 98.6, *vs.Temperature)
	assert.Equal(t, "Oral", vs.TempRoute)

	err = vs.SetTemperature(120.0, "Oral")
	assert.Error(t, err)

	err = vs.SetTemperature(98.6, "Invalid")
	assert.Error(t, err)
}

func TestVitalSign_IsAbnormal(t *testing.T) {
	vs, _ := NewVitalSign(uuid.Must(uuid.NewV7()), uuid.Must(uuid.NewV7()))

	// Normal vitals
	vs.SetBloodPressure(120, 80)
	vs.SetHeartRate(72)
	vs.SetTemperature(98.6, "Oral")
	vs.SetOxygenSat(98)
	vs.SetRespirations(16)
	assert.False(t, vs.IsAbnormal())

	// High BP
	vs2, _ := NewVitalSign(uuid.Must(uuid.NewV7()), uuid.Must(uuid.NewV7()))
	vs2.SetBloodPressure(200, 80)
	assert.True(t, vs2.IsAbnormal())

	// Low O2
	vs3, _ := NewVitalSign(uuid.Must(uuid.NewV7()), uuid.Must(uuid.NewV7()))
	vs3.SetOxygenSat(88)
	assert.True(t, vs3.IsAbnormal())

	// High temp
	vs4, _ := NewVitalSign(uuid.Must(uuid.NewV7()), uuid.Must(uuid.NewV7()))
	vs4.SetTemperature(102.5, "Oral")
	assert.True(t, vs4.IsAbnormal())
}

func TestVitalSign_PainLevel(t *testing.T) {
	vs, _ := NewVitalSign(uuid.Must(uuid.NewV7()), uuid.Must(uuid.NewV7()))

	err := vs.SetPainLevel(5)
	assert.NoError(t, err)
	assert.Equal(t, 5, *vs.PainLevel)

	err = vs.SetPainLevel(11)
	assert.Error(t, err)

	err = vs.SetPainLevel(-1)
	assert.Error(t, err)
}

func TestValidatedVitalSign(t *testing.T) {
	vs, _ := NewVitalSign(uuid.Must(uuid.NewV7()), uuid.Must(uuid.NewV7()))
	vs.SetBloodPressure(120, 80)

	vvs, err := NewValidatedVitalSign(vs)
	require.NoError(t, err)
	assert.Equal(t, vs, vvs.VitalSign())
}
