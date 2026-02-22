package entities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Medication struct {
	Id           uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Name         string
	BrandName    string
	DrugClass    string
	NDCCode      string
	DefaultDose  string
	DefaultRoute string
	Frequency    string
	HighAlert    bool
}

func NewMedication(name, brandName, drugClass, defaultDose, defaultRoute, frequency string, highAlert bool) (*Medication, error) {
	now := time.Now().UTC()
	m := &Medication{
		Id:           uuid.Must(uuid.NewV7()),
		CreatedAt:    now,
		UpdatedAt:    now,
		Name:         name,
		BrandName:    brandName,
		DrugClass:    drugClass,
		DefaultDose:  defaultDose,
		DefaultRoute: defaultRoute,
		Frequency:    frequency,
		HighAlert:    highAlert,
	}
	if err := m.validate(); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Medication) validate() error {
	if m.Name == "" {
		return fmt.Errorf("medication name is required")
	}
	return nil
}

func (m *Medication) UpdateName(name string) error {
	m.Name = name
	m.UpdatedAt = time.Now().UTC()
	return m.validate()
}
