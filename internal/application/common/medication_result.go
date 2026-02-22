package common

type MedicationResult struct {
	Id           string
	Name         string
	BrandName    string
	DrugClass    string
	DefaultDose  string
	DefaultRoute string
	Frequency    string
	HighAlert    bool
}
