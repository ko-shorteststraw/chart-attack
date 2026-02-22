package entities

type ValidatedPatient struct {
	patient *Patient
}

func NewValidatedPatient(p *Patient) (*ValidatedPatient, error) {
	if err := p.validate(); err != nil {
		return nil, err
	}
	return &ValidatedPatient{patient: p}, nil
}

func (v *ValidatedPatient) Patient() *Patient {
	return v.patient
}
