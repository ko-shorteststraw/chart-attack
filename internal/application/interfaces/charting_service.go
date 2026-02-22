package interfaces

import (
	"context"

	"github.com/kendall/chart-attack/internal/application/command"
	"github.com/kendall/chart-attack/internal/application/query"
)

type ChartingService interface {
	RecordVitals(ctx context.Context, cmd command.RecordVitalsCommand) (*command.RecordVitalsResult, error)
	GetPatientVitals(ctx context.Context, q query.GetPatientVitalsQuery) (*query.GetPatientVitalsResult, error)
}
