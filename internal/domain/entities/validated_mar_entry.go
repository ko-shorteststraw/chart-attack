package entities

type ValidatedMAREntry struct {
	marEntry *MAREntry
}

func NewValidatedMAREntry(m *MAREntry) (*ValidatedMAREntry, error) {
	if err := m.validate(); err != nil {
		return nil, err
	}
	return &ValidatedMAREntry{marEntry: m}, nil
}

func (v *ValidatedMAREntry) MAREntry() *MAREntry {
	return v.marEntry
}
