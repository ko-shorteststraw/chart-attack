package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	u, err := NewUser("jsmith", "Jane Smith, RN", "RN", "MedSurg", "BADGE001")
	require.NoError(t, err)
	assert.Equal(t, "jsmith", u.Username)
	assert.Equal(t, "Jane Smith, RN", u.FullName)
	assert.Equal(t, "RN", u.Role)
	assert.True(t, u.Active)
	assert.NotEmpty(t, u.Id)
}

func TestNewUser_ValidationErrors(t *testing.T) {
	_, err := NewUser("", "Jane Smith", "RN", "MedSurg", "")
	assert.Error(t, err)

	_, err = NewUser("jsmith", "", "RN", "MedSurg", "")
	assert.Error(t, err)

	_, err = NewUser("jsmith", "Jane Smith", "InvalidRole", "MedSurg", "")
	assert.Error(t, err)
}

func TestUser_Deactivate(t *testing.T) {
	u, _ := NewUser("jsmith", "Jane Smith", "RN", "MedSurg", "")
	assert.True(t, u.Active)

	u.Deactivate()
	assert.False(t, u.Active)
}
