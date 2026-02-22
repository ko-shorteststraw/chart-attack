package response

import "time"

type VitalSignResponse struct {
	Id             string   `json:"id"`
	CreatedAt      time.Time `json:"createdAt"`
	SystolicBP     *int     `json:"systolicBP"`
	DiastolicBP    *int     `json:"diastolicBP"`
	HeartRate      *int     `json:"heartRate"`
	Temperature    *float64 `json:"temperature"`
	TempRoute      string   `json:"tempRoute"`
	OxygenSat      *int     `json:"oxygenSat"`
	Respirations   *int     `json:"respirations"`
	PainLevel      *int     `json:"painLevel"`
	SupplementalO2 bool     `json:"supplementalO2"`
	O2FlowRate     *float64 `json:"o2FlowRate"`
	Position       string   `json:"position"`
	Notes          string   `json:"notes"`
	IsAbnormal     bool     `json:"isAbnormal"`
}
