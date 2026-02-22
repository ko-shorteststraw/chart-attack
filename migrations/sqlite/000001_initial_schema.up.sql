-- Users (nurses, providers)
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    deleted_at DATETIME,
    username TEXT NOT NULL UNIQUE,
    full_name TEXT NOT NULL,
    role TEXT NOT NULL,
    unit TEXT NOT NULL,
    badge_id TEXT NOT NULL DEFAULT '',
    active INTEGER NOT NULL DEFAULT 1
);

-- Patients
CREATE TABLE IF NOT EXISTS patients (
    id TEXT PRIMARY KEY,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    deleted_at DATETIME,
    mrn TEXT NOT NULL UNIQUE,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    date_of_birth DATETIME NOT NULL,
    room_bed TEXT NOT NULL,
    allergies TEXT NOT NULL DEFAULT '[]',
    diagnoses TEXT NOT NULL DEFAULT '[]',
    code_status TEXT NOT NULL DEFAULT 'Full Code',
    fall_risk INTEGER NOT NULL DEFAULT 0,
    isolation_type TEXT NOT NULL DEFAULT 'None',
    admit_date DATETIME NOT NULL,
    discharge_date DATETIME,
    assigned_nurse_id TEXT REFERENCES users(id)
);

-- Vital Signs
CREATE TABLE IF NOT EXISTS vital_signs (
    id TEXT PRIMARY KEY,
    created_at DATETIME NOT NULL,
    patient_id TEXT NOT NULL REFERENCES patients(id),
    recorded_by TEXT NOT NULL REFERENCES users(id),
    systolic_bp INTEGER,
    diastolic_bp INTEGER,
    heart_rate INTEGER,
    temperature REAL,
    temp_route TEXT,
    oxygen_sat INTEGER,
    respirations INTEGER,
    pain_level INTEGER,
    supplemental_o2 INTEGER NOT NULL DEFAULT 0,
    o2_flow_rate REAL,
    position TEXT,
    notes TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_vital_signs_patient_id ON vital_signs(patient_id);
CREATE INDEX IF NOT EXISTS idx_vital_signs_created_at ON vital_signs(created_at);

-- Audit Entries (append-only)
CREATE TABLE IF NOT EXISTS audit_entries (
    id TEXT PRIMARY KEY,
    created_at DATETIME NOT NULL,
    user_id TEXT NOT NULL,
    patient_id TEXT,
    action TEXT NOT NULL,
    entity_type TEXT NOT NULL,
    entity_id TEXT NOT NULL,
    fields_changed TEXT NOT NULL DEFAULT '{}',
    ip_address TEXT NOT NULL DEFAULT '',
    user_agent TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_audit_entries_patient_id ON audit_entries(patient_id);
CREATE INDEX IF NOT EXISTS idx_audit_entries_entity ON audit_entries(entity_type, entity_id);

-- Idempotency Records
CREATE TABLE IF NOT EXISTS idempotency_records (
    id TEXT PRIMARY KEY,
    created_at DATETIME NOT NULL,
    idempotency_key TEXT NOT NULL UNIQUE,
    response TEXT NOT NULL DEFAULT '{}'
);

CREATE INDEX IF NOT EXISTS idx_idempotency_key ON idempotency_records(idempotency_key);
