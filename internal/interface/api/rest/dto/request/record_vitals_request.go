package request

type RecordVitalsRequest struct {
	SystolicBP     *int     `json:"systolic"`
	DiastolicBP    *int     `json:"diastolic"`
	HeartRate      *int     `json:"hr"`
	Temperature    *float64 `json:"temp"`
	TempRoute      string   `json:"tempRoute"`
	OxygenSat      *int     `json:"o2"`
	Respirations   *int     `json:"resp"`
	PainLevel      *int     `json:"pain"`
	SupplementalO2 bool     `json:"supplementalO2"`
	O2FlowRate     *float64 `json:"o2FlowRate"`
	Position       string   `json:"position"`
	Notes          string   `json:"notes"`
}
