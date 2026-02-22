package repositories

import (
	"context"

	"github.com/kendall/chart-attack/internal/domain/entities"
)

type IdempotencyRepository interface {
	FindByKey(ctx context.Context, key string) (*entities.IdempotencyRecord, error)
	Save(ctx context.Context, record *entities.IdempotencyRecord) error
}
