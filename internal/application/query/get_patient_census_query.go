package query

import "github.com/kendall/chart-attack/internal/application/common"

type GetPatientCensusQuery struct{}

type GetPatientCensusResult struct {
	Patients []common.PatientResult
}
