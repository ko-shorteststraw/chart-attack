package entities

type ValidatedMedication struct {
	medication *Medication
}

func NewValidatedMedication(m *Medication) (*ValidatedMedication, error) {
	if err := m.validate(); err != nil {
		return nil, err
	}
	return &ValidatedMedication{medication: m}, nil
}

func (v *ValidatedMedication) Medication() *Medication {
	return v.medication
}
