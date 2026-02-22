package common

import "time"

type VaccinationResult struct {
	Id               string
	CreatedAt        time.Time
	PatientId        string
	VaccineName      string
	DateAdministered time.Time
	LotNumber        string
	Site             string
	Notes            string
}
