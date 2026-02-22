package entities

type ValidatedVitalSign struct {
	vitalSign *VitalSign
}

func NewValidatedVitalSign(vs *VitalSign) (*ValidatedVitalSign, error) {
	if err := vs.validate(); err != nil {
		return nil, err
	}
	return &ValidatedVitalSign{vitalSign: vs}, nil
}

func (v *ValidatedVitalSign) VitalSign() *VitalSign {
	return v.vitalSign
}
