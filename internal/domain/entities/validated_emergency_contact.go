package entities

type ValidatedEmergencyContact struct {
	contact *EmergencyContact
}

func NewValidatedEmergencyContact(ec *EmergencyContact) (*ValidatedEmergencyContact, error) {
	if err := ec.validate(); err != nil {
		return nil, err
	}
	return &ValidatedEmergencyContact{contact: ec}, nil
}

func (v *ValidatedEmergencyContact) EmergencyContact() *EmergencyContact {
	return v.contact
}
