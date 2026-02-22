package request

type AdmitPatientRequest struct {
	MRN           string   `json:"mrn"`
	FirstName     string   `json:"firstName"`
	LastName      string   `json:"lastName"`
	DateOfBirth   string   `json:"dateOfBirth"`
	RoomBed       string   `json:"roomBed"`
	Allergies     []string `json:"allergies"`
	Diagnoses     []string `json:"diagnoses"`
	CodeStatus    string   `json:"codeStatus"`
	FallRisk      bool     `json:"fallRisk"`
	IsolationType string   `json:"isolationType"`
}
