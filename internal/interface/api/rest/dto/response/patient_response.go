package response

type PatientResponse struct {
	Id            string   `json:"id"`
	MRN           string   `json:"mrn"`
	FirstName     string   `json:"firstName"`
	LastName      string   `json:"lastName"`
	FullName      string   `json:"fullName"`
	RoomBed       string   `json:"roomBed"`
	Allergies     []string `json:"allergies"`
	Diagnoses     []string `json:"diagnoses"`
	CodeStatus    string   `json:"codeStatus"`
	FallRisk      bool     `json:"fallRisk"`
	IsolationType string   `json:"isolationType"`
	IsAdmitted    bool     `json:"isAdmitted"`
}
