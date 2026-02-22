-- name: GetIdempotencyRecord :one
SELECT id, created_at, idempotency_key, response
FROM idempotency_records
WHERE idempotency_key = ?;

-- name: CreateIdempotencyRecord :exec
INSERT INTO idempotency_records (id, created_at, idempotency_key, response)
VALUES (?, ?, ?, ?);
