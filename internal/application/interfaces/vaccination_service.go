package interfaces

import (
	"context"

	"github.com/kendall/chart-attack/internal/application/command"
	"github.com/kendall/chart-attack/internal/application/query"
)

type VaccinationService interface {
	AddVaccination(ctx context.Context, cmd command.AddVaccinationCommand) (*command.AddVaccinationResult, error)
	GetPatientVaccinations(ctx context.Context, q query.GetPatientVaccinationsQuery) (*query.GetPatientVaccinationsResult, error)
}
