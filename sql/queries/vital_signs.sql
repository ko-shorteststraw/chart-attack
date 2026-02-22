-- name: CreateVitalSign :exec
INSERT INTO vital_signs (id, created_at, patient_id, recorded_by, systolic_bp, diastolic_bp, heart_rate, temperature, temp_route, oxygen_sat, respirations, pain_level, supplemental_o2, o2_flow_rate, position, notes)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetVitalSignsByPatientID :many
SELECT id, created_at, patient_id, recorded_by, systolic_bp, diastolic_bp, heart_rate, temperature, temp_route, oxygen_sat, respirations, pain_level, supplemental_o2, o2_flow_rate, position, notes
FROM vital_signs
WHERE patient_id = ?
ORDER BY created_at DESC;

-- name: GetLatestVitalSign :one
SELECT id, created_at, patient_id, recorded_by, systolic_bp, diastolic_bp, heart_rate, temperature, temp_route, oxygen_sat, respirations, pain_level, supplemental_o2, o2_flow_rate, position, notes
FROM vital_signs
WHERE patient_id = ?
ORDER BY created_at DESC
LIMIT 1;

-- name: GetVitalSignByID :one
SELECT id, created_at, patient_id, recorded_by, systolic_bp, diastolic_bp, heart_rate, temperature, temp_route, oxygen_sat, respirations, pain_level, supplemental_o2, o2_flow_rate, position, notes
FROM vital_signs
WHERE id = ?;
