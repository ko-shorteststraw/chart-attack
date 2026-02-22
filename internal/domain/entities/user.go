package entities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

var validRoles = map[string]bool{
	"RN": true, "LPN": true, "CNA": true, "MD": true,
	"PA": true, "NP": true, "Charge": true, "Admin": true,
}

type User struct {
	Id        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string
	FullName  string
	Role      string
	Unit      string
	BadgeId   string
	Active    bool
}

func NewUser(username, fullName, role, unit, badgeId string) (*User, error) {
	now := time.Now().UTC()
	u := &User{
		Id:        uuid.Must(uuid.NewV7()),
		CreatedAt: now,
		UpdatedAt: now,
		Username:  username,
		FullName:  fullName,
		Role:      role,
		Unit:      unit,
		BadgeId:   badgeId,
		Active:    true,
	}
	if err := u.validate(); err != nil {
		return nil, err
	}
	return u, nil
}

func (u *User) validate() error {
	if u.Username == "" {
		return fmt.Errorf("username is required")
	}
	if u.FullName == "" {
		return fmt.Errorf("full name is required")
	}
	if !validRoles[u.Role] {
		return fmt.Errorf("invalid role: %s", u.Role)
	}
	if u.Unit == "" {
		return fmt.Errorf("unit is required")
	}
	return nil
}

func (u *User) UpdateFullName(name string) error {
	u.FullName = name
	u.UpdatedAt = time.Now().UTC()
	return u.validate()
}

func (u *User) UpdateRole(role string) error {
	u.Role = role
	u.UpdatedAt = time.Now().UTC()
	return u.validate()
}

func (u *User) Deactivate() {
	u.Active = false
	u.UpdatedAt = time.Now().UTC()
}
