-- name: CreateUser :exec
INSERT INTO users (id, created_at, updated_at, username, full_name, role, unit, badge_id, active)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetUserByID :one
SELECT id, created_at, updated_at, username, full_name, role, unit, badge_id, active
FROM users
WHERE id = ? AND deleted_at IS NULL;

-- name: GetAllActiveUsers :many
SELECT id, created_at, updated_at, username, full_name, role, unit, badge_id, active
FROM users
WHERE active = 1 AND deleted_at IS NULL
ORDER BY full_name;

-- name: UpdateUser :exec
UPDATE users
SET updated_at = ?, username = ?, full_name = ?, role = ?, unit = ?, badge_id = ?, active = ?
WHERE id = ? AND deleted_at IS NULL;

-- name: SoftDeleteUser :exec
UPDATE users SET deleted_at = ? WHERE id = ?;
