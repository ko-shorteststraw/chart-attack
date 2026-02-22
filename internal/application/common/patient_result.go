package common

import "time"

type PatientResult struct {
	Id            string
	MRN           string
	FirstName     string
	LastName      string
	FullName      string
	DateOfBirth   time.Time
	RoomBed       string
	Allergies     []string
	Diagnoses     []string
	CodeStatus    string
	FallRisk      bool
	IsolationType string
	AdmitDate     time.Time
	IsAdmitted    bool
}
