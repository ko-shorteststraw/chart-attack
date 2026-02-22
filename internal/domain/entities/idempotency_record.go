package entities

import (
	"time"

	"github.com/google/uuid"
)

type IdempotencyRecord struct {
	Id             uuid.UUID
	CreatedAt      time.Time
	IdempotencyKey string
	Response       string
}

func NewIdempotencyRecord(key, response string) *IdempotencyRecord {
	return &IdempotencyRecord{
		Id:             uuid.Must(uuid.NewV7()),
		CreatedAt:      time.Now().UTC(),
		IdempotencyKey: key,
		Response:       response,
	}
}
