-- name: CreateAuditEntry :exec
INSERT INTO audit_entries (id, created_at, user_id, patient_id, action, entity_type, entity_id, fields_changed, ip_address, user_agent)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetAuditEntriesByPatientID :many
SELECT id, created_at, user_id, patient_id, action, entity_type, entity_id, fields_changed, ip_address, user_agent
FROM audit_entries
WHERE patient_id = ?
ORDER BY created_at DESC;

-- name: GetAuditEntriesByEntity :many
SELECT id, created_at, user_id, patient_id, action, entity_type, entity_id, fields_changed, ip_address, user_agent
FROM audit_entries
WHERE entity_type = ? AND entity_id = ?
ORDER BY created_at DESC;
