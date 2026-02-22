package entities

type ValidatedHealthHistoryEntry struct {
	entry *HealthHistoryEntry
}

func NewValidatedHealthHistoryEntry(h *HealthHistoryEntry) (*ValidatedHealthHistoryEntry, error) {
	if err := h.validate(); err != nil {
		return nil, err
	}
	return &ValidatedHealthHistoryEntry{entry: h}, nil
}

func (v *ValidatedHealthHistoryEntry) HealthHistoryEntry() *HealthHistoryEntry {
	return v.entry
}
