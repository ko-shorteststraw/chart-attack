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

type sqlcIdempotencyRepository struct {
	q *sqlc.Queries
}

func NewSqlcIdempotencyRepository(db *sql.DB) repositories.IdempotencyRepository {
	return &sqlcIdempotencyRepository{q: sqlc.New(db)}
}

func (r *sqlcIdempotencyRepository) FindByKey(ctx context.Context, key string) (*entities.IdempotencyRecord, error) {
	row, err := r.q.GetIdempotencyRecord(ctx, key)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("finding idempotency record: %w", err)
	}

	id, err := uuid.Parse(row.ID)
	if err != nil {
		return nil, fmt.Errorf("parsing idempotency record ID: %w", err)
	}
	createdAt, err := time.Parse(time.RFC3339, row.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("parsing created_at: %w", err)
	}

	return &entities.IdempotencyRecord{
		Id:             id,
		CreatedAt:      createdAt,
		IdempotencyKey: row.IdempotencyKey,
		Response:       row.Response,
	}, nil
}

func (r *sqlcIdempotencyRepository) Save(ctx context.Context, record *entities.IdempotencyRecord) error {
	return r.q.CreateIdempotencyRecord(ctx, sqlc.CreateIdempotencyRecordParams{
		ID:             record.Id.String(),
		CreatedAt:      record.CreatedAt.Format(time.RFC3339),
		IdempotencyKey: record.IdempotencyKey,
		Response:       record.Response,
	})
}
