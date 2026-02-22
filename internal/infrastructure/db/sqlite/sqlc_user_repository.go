package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kendall/chart-attack/internal/domain/entities"
	"github.com/kendall/chart-attack/internal/domain/repositories"
	"github.com/kendall/chart-attack/internal/infrastructure/db/sqlc"
)

type sqlcUserRepository struct {
	q *sqlc.Queries
}

func NewSqlcUserRepository(db *sql.DB) repositories.UserRepository {
	return &sqlcUserRepository{q: sqlc.New(db)}
}

func (r *sqlcUserRepository) Save(ctx context.Context, vu *entities.ValidatedUser) error {
	u := vu.User()
	active := int64(0)
	if u.Active {
		active = 1
	}
	return r.q.CreateUser(ctx, sqlc.CreateUserParams{
		ID:        u.Id.String(),
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
		Username:  u.Username,
		FullName:  u.FullName,
		Role:      u.Role,
		Unit:      u.Unit,
		BadgeID:   u.BadgeId,
		Active:    active,
	})
}

func (r *sqlcUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	row, err := r.q.GetUserByID(ctx, id.String())
	if err != nil {
		return nil, fmt.Errorf("finding user by ID: %w", err)
	}
	return mapSqlcUserToEntity(row)
}

func (r *sqlcUserRepository) FindAllActive(ctx context.Context) ([]*entities.User, error) {
	rows, err := r.q.GetAllActiveUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("finding all active users: %w", err)
	}
	users := make([]*entities.User, 0, len(rows))
	for _, row := range rows {
		u, err := mapSqlcUserToEntity(row)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func mapSqlcUserToEntity(row sqlc.User) (*entities.User, error) {
	id, err := uuid.Parse(row.ID)
	if err != nil {
		return nil, fmt.Errorf("parsing user ID: %w", err)
	}
	createdAt, err := time.Parse(time.RFC3339, row.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("parsing created_at: %w", err)
	}
	updatedAt, err := time.Parse(time.RFC3339, row.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("parsing updated_at: %w", err)
	}
	return &entities.User{
		Id:        id,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Username:  row.Username,
		FullName:  row.FullName,
		Role:      row.Role,
		Unit:      row.Unit,
		BadgeId:   row.BadgeID,
		Active:    row.Active == 1,
	}, nil
}
