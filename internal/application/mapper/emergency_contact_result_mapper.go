package mapper

import (
	"github.com/kendall/chart-attack/internal/application/common"
	"github.com/kendall/chart-attack/internal/domain/entities"
)

func ToEmergencyContactResult(ec *entities.EmergencyContact) common.EmergencyContactResult {
	return common.EmergencyContactResult{
		Id:           ec.Id.String(),
		PatientId:    ec.PatientId.String(),
		Name:         ec.Name,
		Relationship: ec.Relationship,
		Phone:        ec.Phone,
		Email:        ec.Email,
		IsPrimary:    ec.IsPrimary,
	}
}

func ToEmergencyContactResults(contacts []*entities.EmergencyContact) []common.EmergencyContactResult {
	results := make([]common.EmergencyContactResult, 0, len(contacts))
	for _, ec := range contacts {
		results = append(results, ToEmergencyContactResult(ec))
	}
	return results
}
