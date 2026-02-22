package common

type EmergencyContactResult struct {
	Id           string
	PatientId    string
	Name         string
	Relationship string
	Phone        string
	Email        string
	IsPrimary    bool
}
