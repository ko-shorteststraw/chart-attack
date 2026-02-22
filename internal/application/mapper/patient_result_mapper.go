package mapper

import (
	"github.com/kendall/chart-attack/internal/application/common"
	"github.com/kendall/chart-attack/internal/domain/entities"
)

func ToPatientResult(p *entities.Patient) common.PatientResult {
	return common.PatientResult{
		Id:            p.Id.String(),
		MRN:           p.MRN,
		FirstName:     p.FirstName,
		LastName:      p.LastName,
		FullName:      p.FullName(),
		DateOfBirth:   p.DateOfBirth,
		RoomBed:       p.RoomBed,
		Allergies:     p.Allergies,
		Diagnoses:     p.Diagnoses,
		CodeStatus:    p.CodeStatus,
		FallRisk:      p.FallRisk,
		IsolationType: p.IsolationType,
		AdmitDate:     p.AdmitDate,
		IsAdmitted:    p.IsAdmitted(),
	}
}

func ToPatientResults(patients []*entities.Patient) []common.PatientResult {
	results := make([]common.PatientResult, 0, len(patients))
	for _, p := range patients {
		results = append(results, ToPatientResult(p))
	}
	return results
}
