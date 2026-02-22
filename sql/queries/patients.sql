-- name: CreatePatient :exec
INSERT INTO patients (id, created_at, updated_at, mrn, first_name, last_name, date_of_birth, room_bed, allergies, diagnoses, code_status, fall_risk, isolation_type, admit_date, discharge_date, assigned_nurse_id)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetPatientByID :one
SELECT id, created_at, updated_at, mrn, first_name, last_name, date_of_birth, room_bed, allergies, diagnoses, code_status, fall_risk, isolation_type, admit_date, discharge_date, assigned_nurse_id
FROM patients
WHERE id = ? AND deleted_at IS NULL;

-- name: GetAdmittedPatients :many
SELECT id, created_at, updated_at, mrn, first_name, last_name, date_of_birth, room_bed, allergies, diagnoses, code_status, fall_risk, isolation_type, admit_date, discharge_date, assigned_nurse_id
FROM patients
WHERE discharge_date IS NULL AND deleted_at IS NULL
ORDER BY room_bed;

-- name: UpdatePatient :exec
UPDATE patients
SET updated_at = ?, first_name = ?, last_name = ?, room_bed = ?, allergies = ?, diagnoses = ?, code_status = ?, fall_risk = ?, isolation_type = ?, assigned_nurse_id = ?
WHERE id = ? AND deleted_at IS NULL;

-- name: DischargePatient :exec
UPDATE patients
SET updated_at = ?, discharge_date = ?
WHERE id = ? AND deleted_at IS NULL;

-- name: SoftDeletePatient :exec
UPDATE patients SET deleted_at = ? WHERE id = ?;
