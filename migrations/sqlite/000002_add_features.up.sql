-- Medications
CREATE TABLE IF NOT EXISTS medications (
    id TEXT PRIMARY KEY,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    deleted_at DATETIME,
    name TEXT NOT NULL,
    brand_name TEXT NOT NULL DEFAULT '',
    drug_class TEXT NOT NULL DEFAULT '',
    ndc_code TEXT NOT NULL DEFAULT '',
    default_dose TEXT NOT NULL DEFAULT '',
    default_route TEXT NOT NULL DEFAULT 'PO',
    frequency TEXT NOT NULL DEFAULT '',
    high_alert INTEGER NOT NULL DEFAULT 0
);

-- MAR Entries
CREATE TABLE IF NOT EXISTS mar_entries (
    id TEXT PRIMARY KEY,
    created_at DATETIME NOT NULL,
    patient_id TEXT NOT NULL REFERENCES patients(id),
    medication_id TEXT NOT NULL REFERENCES medications(id),
    scheduled_time DATETIME NOT NULL,
    administered_at DATETIME,
    administered_by TEXT REFERENCES users(id),
    status TEXT NOT NULL DEFAULT 'Scheduled',
    dose TEXT NOT NULL DEFAULT '',
    route TEXT NOT NULL DEFAULT 'PO',
    site TEXT NOT NULL DEFAULT '',
    hold_reason TEXT NOT NULL DEFAULT '',
    notes TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_mar_entries_patient_id ON mar_entries(patient_id);
CREATE INDEX IF NOT EXISTS idx_mar_entries_scheduled_time ON mar_entries(scheduled_time);

-- Vaccination Records
CREATE TABLE IF NOT EXISTS vaccination_records (
    id TEXT PRIMARY KEY,
    created_at DATETIME NOT NULL,
    patient_id TEXT NOT NULL REFERENCES patients(id),
    recorded_by TEXT NOT NULL REFERENCES users(id),
    vaccine_name TEXT NOT NULL,
    date_administered DATETIME NOT NULL,
    lot_number TEXT NOT NULL DEFAULT '',
    site TEXT NOT NULL DEFAULT '',
    notes TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_vaccination_records_patient_id ON vaccination_records(patient_id);

-- Health History
CREATE TABLE IF NOT EXISTS health_history (
    id TEXT PRIMARY KEY,
    created_at DATETIME NOT NULL,
    patient_id TEXT NOT NULL REFERENCES patients(id),
    recorded_by TEXT NOT NULL REFERENCES users(id),
    condition TEXT NOT NULL,
    date_occurred DATETIME,
    description TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'Active'
);

CREATE INDEX IF NOT EXISTS idx_health_history_patient_id ON health_history(patient_id);

-- Emergency Contacts
CREATE TABLE IF NOT EXISTS emergency_contacts (
    id TEXT PRIMARY KEY,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    patient_id TEXT NOT NULL REFERENCES patients(id),
    name TEXT NOT NULL,
    relationship TEXT NOT NULL DEFAULT '',
    phone TEXT NOT NULL DEFAULT '',
    email TEXT NOT NULL DEFAULT '',
    is_primary INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_emergency_contacts_patient_id ON emergency_contacts(patient_id);

-- Tasks
CREATE TABLE IF NOT EXISTS tasks (
    id TEXT PRIMARY KEY,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    patient_id TEXT NOT NULL REFERENCES patients(id),
    assigned_to TEXT NOT NULL REFERENCES users(id),
    title TEXT NOT NULL,
    category TEXT NOT NULL DEFAULT 'Custom',
    due_at DATETIME NOT NULL,
    completed_at DATETIME,
    completed_by TEXT REFERENCES users(id),
    priority TEXT NOT NULL DEFAULT 'Routine',
    recurring INTEGER NOT NULL DEFAULT 0,
    recur_interval TEXT,
    notes TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_tasks_patient_id ON tasks(patient_id);
CREATE INDEX IF NOT EXISTS idx_tasks_assigned_to ON tasks(assigned_to);
CREATE INDEX IF NOT EXISTS idx_tasks_due_at ON tasks(due_at);
