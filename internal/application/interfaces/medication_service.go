package interfaces

import (
	"context"

	"github.com/kendall/chart-attack/internal/application/command"
	"github.com/kendall/chart-attack/internal/application/query"
)

type MedicationService interface {
	AddMedication(ctx context.Context, cmd command.AddMedicationCommand) (*command.AddMedicationResult, error)
	AdministerMedication(ctx context.Context, cmd command.AdministerMedicationCommand) (*command.AdministerMedicationResult, error)
	GetPatientMAR(ctx context.Context, q query.GetPatientMARQuery) (*query.GetPatientMARResult, error)
	GetAllMedications(ctx context.Context, q query.GetAllMedicationsQuery) (*query.GetAllMedicationsResult, error)
}
