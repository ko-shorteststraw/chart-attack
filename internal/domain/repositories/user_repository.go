package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/kendall/chart-attack/internal/domain/entities"
)

type UserRepository interface {
	Save(ctx context.Context, user *entities.ValidatedUser) error
	FindByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	FindAllActive(ctx context.Context) ([]*entities.User, error)
}
