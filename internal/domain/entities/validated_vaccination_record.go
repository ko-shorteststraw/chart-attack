package entities

type ValidatedVaccinationRecord struct {
	vaccinationRecord *VaccinationRecord
}

func NewValidatedVaccinationRecord(v *VaccinationRecord) (*ValidatedVaccinationRecord, error) {
	if err := v.validate(); err != nil {
		return nil, err
	}
	return &ValidatedVaccinationRecord{vaccinationRecord: v}, nil
}

func (vv *ValidatedVaccinationRecord) VaccinationRecord() *VaccinationRecord {
	return vv.vaccinationRecord
}
