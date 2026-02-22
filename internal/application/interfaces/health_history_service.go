package interfaces

import (
	"context"

	"github.com/kendall/chart-attack/internal/application/command"
	"github.com/kendall/chart-attack/internal/application/query"
)

type HealthHistoryService interface {
	AddHealthHistory(ctx context.Context, cmd command.AddHealthHistoryCommand) (*command.AddHealthHistoryResult, error)
	GetPatientHealthHistory(ctx context.Context, q query.GetPatientHealthHistoryQuery) (*query.GetPatientHealthHistoryResult, error)
}
