package command

import "github.com/google/uuid"

type AddMedicationCommand struct {
	IdempotencyKey string
	Name           string
	BrandName      string
	DrugClass      string
	DefaultDose    string
	DefaultRoute   string
	Frequency      string
	HighAlert      bool
	AddedBy        uuid.UUID
}

type AddMedicationResult struct {
	MedicationId string
}
