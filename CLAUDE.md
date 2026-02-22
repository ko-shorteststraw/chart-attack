# Chart Attack - Nursing Charting Application

## Project Overview

Chart Attack is a nursing charting and clinical documentation application designed for bedside nurses. It prioritizes speed, safety, and hands-free operation during patient care. The app handles vital signs, medication administration, assessments, shift handoffs, and secure team communication.

---

## Tech Stack

| Layer | Technology | Notes |
|-------|-----------|-------|
| **Language** | Go 1.24+ | Backend and templating |
| **Architecture** | DDD / Onion / CQRS | Based on [go-ddd](https://github.com/sklinkert/go-ddd) |
| **HTTP Framework** | Echo v4 | Router, middleware, static files |
| **Database (primary)** | SQLite | Via `modernc.org/sqlite` (pure Go, no CGO) |
| **Database (future)** | PostgreSQL | Via `pgx/v5`; repository interfaces abstract the DB |
| **SQL Codegen** | sqlc | Type-safe SQL → Go code generation |
| **Migrations** | golang-migrate | `migrations/` directory, numbered files |
| **Frontend Reactivity** | Datastar v1 | Hypermedia framework, ~11KB, SSE-driven |
| **Go SDK for Datastar** | `github.com/starfederation/datastar-go` | SSE helpers, signal reading, element patching |
| **Templating** | Go `html/template` or `templ` | Server-rendered HTML fragments |
| **CSS** | Pico CSS v2 | Classless/semantic CSS, dark mode built-in |
| **IDs** | UUIDv7 | `github.com/google/uuid`, time-sortable |
| **Testing** | `testify`, `testcontainers-go` (for Postgres) | Unit + integration |

---

## Architecture: DDD Onion Model

Follow the [go-ddd](https://github.com/sklinkert/go-ddd) reference architecture strictly. Dependencies point inward only.

```
┌─────────────────────────────────────────────┐
│  Interface Layer (HTTP controllers, SSE)    │
│  ┌───────────────────────────────────────┐  │
│  │  Application Layer (services, CQRS)   │  │
│  │  ┌─────────────────────────────────┐  │  │
│  │  │  Domain Layer (entities, repos) │  │  │
│  │  └─────────────────────────────────┘  │  │
│  └───────────────────────────────────────┘  │
│  Infrastructure Layer (DB, external APIs)   │
└─────────────────────────────────────────────┘
```

### Layer Rules

- **Domain** (innermost): Zero imports from other project layers. Pure business logic, entities, value objects, and repository interfaces.
- **Application**: Depends on Domain only. Contains services, commands, queries, DTOs, and mappers.
- **Infrastructure**: Implements domain repository interfaces. Depends on Domain layer for interfaces.
- **Interface**: Depends on Application (via service interfaces). HTTP controllers, SSE handlers, request/response DTOs.

---

## Directory Structure

```
chart-attack/
├── CLAUDE.md
├── cmd/
│   └── server/
│       └── main.go                          # Entrypoint, manual DI wiring
├── internal/
│   ├── domain/
│   │   ├── entities/
│   │   │   ├── patient.go                   # Patient aggregate
│   │   │   ├── validated_patient.go
│   │   │   ├── vital_sign.go                # BP, HR, temp, O2, resp
│   │   │   ├── validated_vital_sign.go
│   │   │   ├── medication.go                # Medication definition
│   │   │   ├── mar_entry.go                 # Medication administration record entry
│   │   │   ├── validated_mar_entry.go
│   │   │   ├── assessment.go                # Head-to-toe assessment
│   │   │   ├── validated_assessment.go
│   │   │   ├── pain_assessment.go           # 0-10 scale, Wong-Baker
│   │   │   ├── validated_pain_assessment.go
│   │   │   ├── intake_output.go             # I&O tracking
│   │   │   ├── validated_intake_output.go
│   │   │   ├── wound.go                     # Wound documentation
│   │   │   ├── validated_wound.go
│   │   │   ├── task.go                      # Patient task/checklist item
│   │   │   ├── validated_task.go
│   │   │   ├── shift_report.go              # SBAR handoff report
│   │   │   ├── validated_shift_report.go
│   │   │   ├── message.go                   # Secure team message
│   │   │   ├── validated_message.go
│   │   │   ├── notification.go              # Physician notification log
│   │   │   ├── validated_notification.go
│   │   │   ├── family_note.go               # Family communication note
│   │   │   ├── validated_family_note.go
│   │   │   ├── user.go                      # Nurse/provider user
│   │   │   ├── validated_user.go
│   │   │   ├── audit_entry.go               # Audit trail entry
│   │   │   └── idempotency_record.go
│   │   └── repositories/
│   │       ├── patient_repository.go
│   │       ├── vital_sign_repository.go
│   │       ├── mar_repository.go
│   │       ├── assessment_repository.go
│   │       ├── pain_assessment_repository.go
│   │       ├── intake_output_repository.go
│   │       ├── wound_repository.go
│   │       ├── task_repository.go
│   │       ├── shift_report_repository.go
│   │       ├── message_repository.go
│   │       ├── notification_repository.go
│   │       ├── family_note_repository.go
│   │       ├── user_repository.go
│   │       ├── audit_repository.go
│   │       └── idempotency_repository.go
│   ├── application/
│   │   ├── command/                         # Write operations
│   │   │   ├── record_vitals_command.go
│   │   │   ├── administer_medication_command.go
│   │   │   ├── create_assessment_command.go
│   │   │   ├── record_pain_command.go
│   │   │   ├── record_intake_output_command.go
│   │   │   ├── document_wound_command.go
│   │   │   ├── create_task_command.go
│   │   │   ├── complete_task_command.go
│   │   │   ├── create_shift_report_command.go
│   │   │   ├── send_message_command.go
│   │   │   ├── log_notification_command.go
│   │   │   ├── admit_patient_command.go
│   │   │   ├── discharge_patient_command.go
│   │   │   └── ...
│   │   ├── query/                           # Read operations
│   │   │   ├── get_patient_census_query.go
│   │   │   ├── get_patient_vitals_query.go
│   │   │   ├── get_patient_mar_query.go
│   │   │   ├── get_shift_summary_query.go
│   │   │   ├── get_abnormal_values_query.go
│   │   │   ├── get_task_checklist_query.go
│   │   │   ├── get_messages_query.go
│   │   │   └── ...
│   │   ├── common/                          # Shared result DTOs
│   │   │   ├── patient_result.go
│   │   │   ├── vital_sign_result.go
│   │   │   └── ...
│   │   ├── mapper/                          # Entity → result DTO mappers
│   │   │   ├── patient_result_mapper.go
│   │   │   └── ...
│   │   ├── interfaces/                      # Service interfaces
│   │   │   ├── charting_service.go
│   │   │   ├── patient_service.go
│   │   │   ├── medication_service.go
│   │   │   ├── task_service.go
│   │   │   ├── communication_service.go
│   │   │   ├── reporting_service.go
│   │   │   └── drug_interaction_service.go
│   │   └── services/                        # Service implementations
│   │       ├── charting_service.go
│   │       ├── patient_service.go
│   │       ├── medication_service.go
│   │       ├── task_service.go
│   │       ├── communication_service.go
│   │       ├── reporting_service.go
│   │       └── drug_interaction_service.go
│   ├── infrastructure/
│   │   ├── db/
│   │   │   ├── sqlite/
│   │   │   │   ├── connection.go
│   │   │   │   ├── sqlc_patient_repository.go
│   │   │   │   ├── sqlc_vital_sign_repository.go
│   │   │   │   └── ...
│   │   │   ├── postgres/                    # Future: PostgreSQL implementations
│   │   │   │   └── ...
│   │   │   └── sqlc/                        # Auto-generated by sqlc
│   │   │       ├── db.go
│   │   │       ├── models.go
│   │   │       └── *.sql.go
│   │   ├── storage/
│   │   │   └── photo_storage.go             # Wound photo file storage
│   │   └── external/
│   │       └── drug_interaction_client.go   # External drug interaction API
│   └── interface/
│       └── api/
│           ├── rest/
│           │   ├── patient_controller.go
│           │   ├── charting_controller.go
│           │   ├── medication_controller.go
│           │   ├── task_controller.go
│           │   ├── communication_controller.go
│           │   ├── reporting_controller.go
│           │   └── dto/
│           │       ├── request/
│           │       ├── response/
│           │       └── mapper/
│           └── sse/
│               ├── vital_sign_handler.go    # Real-time vital sign updates
│               ├── task_handler.go          # Live task list updates
│               ├── message_handler.go       # Real-time messaging
│               └── notification_handler.go  # Timed reminders, alerts
├── templates/
│   ├── layouts/
│   │   └── base.html                        # Base layout with Pico CSS + Datastar
│   ├── pages/
│   │   ├── dashboard.html                   # Patient census / room overview
│   │   ├── patient.html                     # Single patient chart view
│   │   ├── mar.html                         # Medication administration record
│   │   ├── assessment.html                  # Head-to-toe assessment form
│   │   ├── vitals.html                      # Vital signs entry/history
│   │   ├── tasks.html                       # Task checklist
│   │   ├── handoff.html                     # Shift handoff / SBAR
│   │   ├── messages.html                    # Secure messaging
│   │   └── reports.html                     # Shift summary / exports
│   └── partials/
│       ├── vital_sign_row.html              # SSE-patchable vital sign row
│       ├── task_item.html                   # SSE-patchable task item
│       ├── patient_card.html                # Census card with flags
│       ├── pain_scale.html                  # Pain assessment widget
│       ├── drug_alert.html                  # Drug interaction alert banner
│       ├── message_bubble.html              # Chat message
│       └── ...
├── static/
│   ├── css/
│   │   └── app.css                          # Custom overrides on top of Pico
│   ├── js/
│   │   └── datastar.min.js                  # Datastar CDN or vendored
│   └── images/
│       └── wong-baker/                      # Wong-Baker faces scale images
├── migrations/
│   ├── sqlite/
│   │   ├── 000001_initial_schema.up.sql
│   │   └── 000001_initial_schema.down.sql
│   └── postgres/
│       ├── 000001_initial_schema.up.sql
│       └── 000001_initial_schema.down.sql
├── sql/
│   └── queries/
│       ├── patients.sql
│       ├── vital_signs.sql
│       ├── mar.sql
│       ├── assessments.sql
│       ├── tasks.sql
│       ├── messages.sql
│       └── ...
├── sqlc.yaml
├── go.mod
└── go.sum
```

---

## Domain Entities

### Core Patterns (from go-ddd)

1. **Factory methods**: `NewX(...)` creates entities with UUIDv7 ID and timestamps set in the domain layer (never the DB).
2. **Private validation**: `validate()` is private, called by mutation methods and `NewValidatedX`.
3. **Validated wrappers**: `ValidatedX` wraps an entity and proves it passed validation. Repository write methods accept only `*ValidatedX`. Read methods return `*X` (no re-validation of historical data).
4. **Mutation methods**: Named `UpdateX(value)` methods that set the field, bump `UpdatedAt`, and call `validate()`.
5. **Idempotency**: All write commands carry an `IdempotencyKey`. Services check before processing and cache results.
6. **Soft deletion**: `deleted_at` column at DB layer only. Domain entities have no `DeletedAt` field.

### Entity Definitions

#### Patient (Aggregate Root)
```go
type Patient struct {
    Id              uuid.UUID
    CreatedAt       time.Time
    UpdatedAt       time.Time
    MRN             string      // Medical record number
    FirstName       string
    LastName        string
    DateOfBirth     time.Time
    RoomBed         string      // e.g. "401-A"
    Allergies       []string    // Allergy list
    Diagnoses       []string    // Active diagnoses
    CodeStatus      string      // "Full Code", "DNR", "DNI", "DNR/DNI", "Comfort Care"
    FallRisk        bool
    IsolationType   string      // "None", "Contact", "Droplet", "Airborne", "Protective"
    AdmitDate       time.Time
    DischargeDate   *time.Time  // nil if still admitted
    AssignedNurseId *uuid.UUID
}
```

#### VitalSign
```go
type VitalSign struct {
    Id              uuid.UUID
    CreatedAt       time.Time
    PatientId       uuid.UUID
    RecordedBy      uuid.UUID   // User ID of nurse
    SystolicBP      *int        // mmHg (nullable — not all vitals taken every time)
    DiastolicBP     *int
    HeartRate       *int        // bpm
    Temperature     *float64    // Fahrenheit
    TempRoute       string      // "Oral", "Tympanic", "Axillary", "Rectal", "Temporal"
    OxygenSat       *int        // SpO2 percentage
    Respirations    *int        // breaths per minute
    PainLevel       *int        // 0-10
    Supplemental02  bool        // On supplemental oxygen
    O2FlowRate      *float64    // L/min if on O2
    Position        string      // "Sitting", "Standing", "Supine", "Left lateral"
    Notes           string
}
```

#### MAREntry (Medication Administration Record)
```go
type MAREntry struct {
    Id              uuid.UUID
    CreatedAt       time.Time
    PatientId       uuid.UUID
    MedicationId    uuid.UUID
    ScheduledTime   time.Time
    AdministeredAt  *time.Time  // nil if not yet given
    AdministeredBy  *uuid.UUID
    Status          string      // "Scheduled", "Given", "Held", "Refused", "Omitted"
    Dose            string      // e.g. "500mg"
    Route           string      // "PO", "IV", "IM", "SQ", "SL", "PR", "Topical", "Inhaled"
    Site            string      // Injection site if applicable
    HoldReason      string      // If held/omitted, reason
    Notes           string
}
```

#### Medication
```go
type Medication struct {
    Id              uuid.UUID
    CreatedAt       time.Time
    UpdatedAt       time.Time
    Name            string      // Generic name
    BrandName       string
    DrugClass       string
    NDCCode         string      // National Drug Code
    DefaultDose     string
    DefaultRoute    string
    Frequency       string      // "BID", "TID", "Q4H", "PRN", etc.
    HighAlertDrug   bool
    Interactions    []string    // Known drug interaction codes/names
}
```

#### Assessment (Head-to-Toe)
```go
type Assessment struct {
    Id              uuid.UUID
    CreatedAt       time.Time
    PatientId       uuid.UUID
    AssessedBy      uuid.UUID
    ShiftType       string      // "Day", "Night", "Evening"

    // Neurological
    NeuroLOC        string      // Level of consciousness: "Alert", "Confused", "Lethargic", "Obtunded", "Comatose"
    NeuroOrientation string     // "Oriented x4", "Oriented x3", etc.
    NeuroPupils     string      // "PERRLA", "Unequal", "Fixed"
    NeuroNotes      string

    // Cardiovascular
    CardiacRhythm   string      // "NSR", "AFib", "SVT", etc.
    CardiacEdema    string      // "None", "Trace", "1+", "2+", "3+", "4+"
    CardiacNotes    string

    // Respiratory
    RespBreathSounds string     // "Clear bilateral", "Diminished bases", "Wheezes", "Crackles"
    RespEffort       string     // "Unlabored", "Labored", "Accessory muscle use"
    RespNotes        string

    // GI
    GIBowelSounds    string     // "Active", "Hypoactive", "Hyperactive", "Absent"
    GIAbdomen        string     // "Soft", "Firm", "Distended", "Tender"
    GILastBM         *time.Time
    GINotes          string

    // GU
    GUOutput         string     // "Voiding", "Foley", "Incontinent"
    GUNotes          string

    // Skin
    SkinIntegrity    string     // "Intact", "Impaired"
    SkinColor        string     // "WNL", "Pale", "Jaundiced", "Cyanotic", "Flushed"
    SkinTurgor       string     // "Good", "Poor", "Tenting"
    SkinNotes        string

    // Musculoskeletal
    MSKMobility      string     // "Independent", "Assist x1", "Assist x2", "Bedrest", "Wheelchair"
    MSKNotes         string

    // Psychosocial
    PsychMood        string     // "Appropriate", "Anxious", "Depressed", "Agitated", "Flat"
    PsychNotes       string

    // IV/Lines
    IVAccess         string     // "PIV", "Central line", "PICC", "Port", "None"
    IVSite           string
    IVSiteCondition  string     // "No redness", "Redness", "Swelling", "Phlebitis"
    IVNotes          string

    Notes            string     // General assessment notes
}
```

#### PainAssessment
```go
type PainAssessment struct {
    Id              uuid.UUID
    CreatedAt       time.Time
    PatientId       uuid.UUID
    AssessedBy      uuid.UUID
    ScaleType       string      // "Numeric", "WongBaker", "FLACC", "CPOT"
    PainLevel       int         // 0-10
    WongBakerFace   *int        // 0-5 face index (if WongBaker scale)
    Location        string      // Body location of pain
    Quality         string      // "Sharp", "Dull", "Aching", "Burning", "Stabbing", "Throbbing"
    Radiation       string      // Does it radiate? Where?
    Duration        string      // "Constant", "Intermittent", onset time
    Aggravating     string      // What makes it worse
    Alleviating     string      // What makes it better
    InterventionId  *uuid.UUID  // Link to MAR entry if med given
    ReassessAt      *time.Time  // When to reassess post-intervention
    Notes           string
}
```

#### IntakeOutput
```go
type IntakeOutput struct {
    Id              uuid.UUID
    CreatedAt       time.Time
    PatientId       uuid.UUID
    RecordedBy      uuid.UUID
    Type            string      // "Intake" or "Output"
    Category        string      // Intake: "PO", "IV", "TubeFeeding", "Blood" / Output: "Urine", "Emesis", "Stool", "Drain", "Blood loss"
    AmountML        int
    Notes           string
}
```

#### Wound
```go
type Wound struct {
    Id              uuid.UUID
    CreatedAt       time.Time
    UpdatedAt       time.Time
    PatientId       uuid.UUID
    DocumentedBy    uuid.UUID
    Location        string      // Body location
    Type            string      // "Pressure injury", "Surgical", "Laceration", "Skin tear", "Diabetic ulcer"
    Stage           string      // For pressure injuries: "Stage 1-4", "Unstageable", "DTI"
    LengthCM        float64
    WidthCM         float64
    DepthCM         float64
    Drainage        string      // "None", "Serous", "Sanguineous", "Serosanguineous", "Purulent"
    DrainageAmount  string      // "Scant", "Small", "Moderate", "Large"
    WoundBed        string      // "Granulation", "Slough", "Eschar", "Epithelial", "Mixed"
    PeriwoundSkin   string      // "Intact", "Macerated", "Erythema", "Indurated"
    Treatment       string      // Dressing/treatment applied
    PhotoPath       string      // File path to wound photo
    Notes           string
}
```

#### Task
```go
type Task struct {
    Id              uuid.UUID
    CreatedAt       time.Time
    UpdatedAt       time.Time
    PatientId       uuid.UUID
    AssignedTo      uuid.UUID
    Title           string
    Category        string      // "Medication", "Vital", "Turn", "Assessment", "Lab", "Custom"
    DueAt           time.Time
    CompletedAt     *time.Time
    CompletedBy     *uuid.UUID
    Priority        string      // "Routine", "Urgent", "STAT"
    Recurring       bool
    RecurInterval   *string     // "Q2H", "Q4H", "Q8H", "BID", "TID", etc.
    Notes           string
}
```

#### ShiftReport (SBAR)
```go
type ShiftReport struct {
    Id              uuid.UUID
    CreatedAt       time.Time
    PatientId       uuid.UUID
    FromNurseId     uuid.UUID
    ToNurseId       uuid.UUID
    ShiftType       string      // "Day→Night", "Night→Day", "Day→Evening", etc.

    // SBAR fields
    Situation       string      // Current status, why patient is here
    Background      string      // Relevant history, recent events
    Assessment      string      // Current assessment findings
    Recommendation  string      // Plan, pending orders, things to watch

    // Quick reference
    CodeStatus      string
    Allergies       string
    ActiveIVs       string
    PendingLabs     string
    DietStatus      string
    MobilityStatus  string
    FallRisk        bool
    IsolationType   string
    Notes           string
}
```

#### Message (Secure Team Messaging)
```go
type Message struct {
    Id              uuid.UUID
    CreatedAt       time.Time
    FromUserId      uuid.UUID
    ToUserId        *uuid.UUID  // nil = group/unit message
    ThreadId        *uuid.UUID  // For threaded conversations
    PatientId       *uuid.UUID  // Optional patient context
    Content         string
    Priority        string      // "Normal", "Urgent"
    ReadAt          *time.Time
}
```

#### User (Nurse/Provider)
```go
type User struct {
    Id              uuid.UUID
    CreatedAt       time.Time
    UpdatedAt       time.Time
    Username        string
    FullName        string
    Role            string      // "RN", "LPN", "CNA", "MD", "PA", "NP", "Charge", "Admin"
    Unit            string      // "ICU", "MedSurg", "ER", "Tele", "L&D", etc.
    BadgeId         string      // For barcode scanner ID verification
    Active          bool
}
```

#### AuditEntry
```go
type AuditEntry struct {
    Id              uuid.UUID
    CreatedAt       time.Time    // Immutable timestamp
    UserId          uuid.UUID
    PatientId       *uuid.UUID
    Action          string       // "CREATE", "UPDATE", "DELETE", "VIEW", "LOGIN", "LOGOUT"
    EntityType      string       // "VitalSign", "MAREntry", "Assessment", etc.
    EntityId        uuid.UUID
    FieldsChanged   string       // JSON of changed fields
    IPAddress       string
    UserAgent       string
}
```

---

## Datastar Frontend Patterns

### Base Layout

Every page includes Pico CSS and Datastar:

```html
<!DOCTYPE html>
<html lang="en" data-theme="light">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="stylesheet" href="/static/css/pico.min.css">
    <link rel="stylesheet" href="/static/css/app.css">
    <title>Chart Attack</title>
</head>
<body>
    <main class="container">
        {{ template "content" . }}
    </main>
    <script src="/static/js/datastar.min.js"></script>
</body>
</html>
```

### Reactivity Pattern: SSE for Live Updates

Use Datastar SSE for all real-time features. The Go backend streams updates via `datastar.NewSSE()`.

```html
<!-- Vital signs table that auto-updates -->
<div id="vitals-table"
     data-on-load="@get('/api/patients/{{.PatientId}}/vitals/stream')">
    <!-- Server patches rows here via SSE -->
</div>
```

```go
// Backend SSE handler
func (c *ChartingController) StreamVitals(w http.ResponseWriter, r *http.Request) {
    sse := datastar.NewSSE(w, r)
    // Patch the vitals table with new HTML when data changes
    sse.PatchElements(renderVitalSignRows(vitals))
}
```

### Forms: Signal Binding + POST

```html
<!-- Vital signs entry form -->
<form data-signals="{systolic: '', diastolic: '', hr: '', temp: '', o2: '', resp: ''}"
      data-on-submit__prevent="@post('/api/patients/{{.PatientId}}/vitals')">
    <input type="number" data-bind:systolic placeholder="Systolic" />
    <input type="number" data-bind:diastolic placeholder="Diastolic" />
    <input type="number" data-bind:hr placeholder="Heart Rate" />
    <button type="submit" data-indicator:saving data-attr:disabled="$saving">
        Record Vitals
    </button>
</form>
```

### Timed Reminders

```html
<!-- Poll for due tasks every 30 seconds -->
<div id="task-alerts"
     data-on-interval.30000ms="@get('/api/tasks/due')">
</div>
```

### Drug Interaction Alerts

```html
<!-- Check interactions on medication selection -->
<select data-bind:medicationId
        data-on-change="@get('/api/medications/interactions', {filterSignals: 'medicationId patientId'})">
    {{range .Medications}}
    <option value="{{.Id}}">{{.Name}}</option>
    {{end}}
</select>
<div id="interaction-alerts">
    <!-- Server patches alerts here -->
</div>
```

### Loading States

```html
<button data-on-click="@post('/api/patients/{{.Id}}/assess')"
        data-indicator:assessing
        data-attr:aria-busy="$assessing"
        data-attr:disabled="$assessing">
    Save Assessment
</button>
```

### Dark Mode Toggle

Pico CSS supports dark mode via `data-theme` attribute on `<html>`:

```html
<html data-signals="{_theme: 'light'}"
      data-attr:data-theme="$_theme">
<!-- Underscore prefix: _theme is local-only, not sent to server -->
<button data-on-click="$_theme = $_theme === 'light' ? 'dark' : 'light'">
    Toggle Dark Mode
</button>
```

---

## Database Schema Conventions

### SQLite Primary / PostgreSQL Abstraction

- Write migrations for **both** SQLite and PostgreSQL in `migrations/sqlite/` and `migrations/postgres/`.
- Use `TEXT` for UUIDs in SQLite, `UUID` type in PostgreSQL.
- Use `INTEGER` for booleans in SQLite (0/1), `BOOLEAN` in PostgreSQL.
- Use `DATETIME` in SQLite, `TIMESTAMPTZ` in PostgreSQL.
- Every table has: `id`, `created_at`, `updated_at` (where mutable), `deleted_at` (soft delete).
- sqlc config in `sqlc.yaml` should have separate configurations per database engine.

### Key Tables

```sql
-- patients
CREATE TABLE patients (
    id TEXT PRIMARY KEY,            -- UUIDv7
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    deleted_at DATETIME,
    mrn TEXT NOT NULL UNIQUE,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    date_of_birth DATETIME NOT NULL,
    room_bed TEXT NOT NULL,
    allergies TEXT NOT NULL DEFAULT '[]',        -- JSON array
    diagnoses TEXT NOT NULL DEFAULT '[]',        -- JSON array
    code_status TEXT NOT NULL DEFAULT 'Full Code',
    fall_risk INTEGER NOT NULL DEFAULT 0,
    isolation_type TEXT NOT NULL DEFAULT 'None',
    admit_date DATETIME NOT NULL,
    discharge_date DATETIME,
    assigned_nurse_id TEXT
);

-- vital_signs
CREATE TABLE vital_signs (
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

-- audit_entries (append-only, no update/delete)
CREATE TABLE audit_entries (
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
```

---

## Feature Implementation Notes

### Core Charting
- Vital signs entry: Simple form → POST → SSE patches the vitals history table.
- MAR: Calendar-style view with time slots. Clicking a slot opens a quick-administer dialog (native `<dialog>`).
- Pain assessment: Interactive 0-10 slider or clickable Wong-Baker faces. Data sent via Datastar signals.
- Head-to-toe assessment: Tabbed or accordion form (using `<details>`/`<summary>`), one section per body system.
- Wound documentation: Form + file input for photo. Photo uploaded separately, path stored on wound record.

### Workflow
- Shift handoff: Pre-populated SBAR template pulling latest vitals, active meds, recent assessments. Nurse edits and saves.
- Task checklists: Server generates recurring tasks (Q2H turns, scheduled vitals). SSE pushes due tasks to the nurse's view. Completing a task is a single button click → POST.
- Timed reminders: `data-on-interval` polls for overdue/upcoming tasks. Browser Notification API for alerts.
- Quick-add shortcuts: Configurable macro buttons (e.g., "Vitals WNL", "Neuro intact") that pre-fill forms with common values.
- Voice-to-text: Use the Web Speech API (`SpeechRecognition`) in a small JS module. Transcribed text fills the notes field via Datastar signals.

### Patient Management
- Census dashboard: Grid of patient cards showing room, name, code status, fall risk, isolation — all color-coded. SSE keeps it live.
- Flags: Code status, fall risk, isolation, and allergies are prominently displayed as colored badges on patient cards using Pico's role attributes or small custom CSS.

### Safety
- Drug interaction alerts: On medication selection, backend checks against a local interaction database or external API. Alerts returned as SSE-patched warning banners.
- Missed charting: Background job (goroutine) checks for overdue tasks/charting. Pushes notifications via SSE.
- Audit trail: Every write operation creates an `AuditEntry` in the service layer. Audit entries are append-only (no update/delete).
- Offline mode: Service Worker caches the app shell and queues failed POSTs in IndexedDB. On reconnect, replays queue. Idempotency keys prevent duplicates.

### Communication
- Secure messaging: Real-time via SSE. Messages stored server-side. No client-side message storage.
- Physician notification logging: Structured form — who was called, when, what was communicated, response received.
- Family notes: Free-text notes tagged with date, author, and patient.

### Reporting
- Shift summary: Auto-generated from the day's charting data — vitals trends, meds given, assessments completed, I&O totals, significant events.
- Abnormal value flagging: Service layer checks vitals against configurable normal ranges. Flagged values get a CSS class (`aria-invalid`) for visual highlighting.
- Export: Generate HTML → PDF (server-side) or CSV for structured data.

---

## Coding Conventions

### Go
- Follow standard Go project layout with `cmd/`, `internal/`.
- No DI container — manual wiring in `main.go`.
- Use `context.Context` throughout for cancellation.
- Errors are returned, not panicked. Wrap with `fmt.Errorf("operation: %w", err)`.
- Log with `slog` (structured logging, stdlib).
- Environment config via env vars or a `.env` file (no viper).

### Naming (from go-ddd)
| Concept | Pattern | Example |
|---------|---------|---------|
| Entity factory | `NewX(...)` | `NewPatient(mrn, name, ...)` |
| Validated wrapper | `NewValidatedX(entity)` | `NewValidatedPatient(patient)` |
| Mutation method | `entity.UpdateX(value)` | `patient.UpdateCodeStatus("DNR")` |
| Command | `VerbEntityCommand` | `RecordVitalsCommand` |
| Query | `GetEntityByFieldQuery` | `GetPatientVitalsQuery` |
| Repository interface | `XRepository` | `VitalSignRepository` |
| Repository impl | `SqlcXRepository` | `SqlcVitalSignRepository` |
| Service interface | `XService` | `ChartingService` |
| Controller | `XController` | `ChartingController` |
| Request DTO | `VerbEntityRequest` | `RecordVitalsRequest` |
| Response DTO | `EntityResponse` | `VitalSignResponse` |

### Testing
- Domain logic: Unit tests with `testify`, no mocks needed (pure functions).
- Services: Unit tests with mock repositories (interfaces make this trivial).
- Repositories: Integration tests with real SQLite (fast, no container needed).
- Controllers: HTTP integration tests using `httptest`.
- Name test files `*_test.go` adjacent to the code they test.

### Frontend
- Prefer server-rendered HTML with Datastar attributes over client-side JS.
- Use Pico CSS semantic elements: `<article>` for cards, `<dialog>` for modals, `<details>` for accordions, `<nav>` for navigation.
- Custom CSS goes in `static/css/app.css` — keep it minimal, only override Pico when necessary.
- All IDs on SSE-patchable elements must be stable and predictable (e.g., `vitals-row-{{.Id}}`).
- Use `_` prefix for local-only signals that shouldn't be sent to the server.

---

## API Routes

```
# Pages (HTML)
GET  /                              → Dashboard (patient census)
GET  /patients/:id                  → Patient chart view
GET  /patients/:id/vitals           → Vitals page
GET  /patients/:id/mar              → MAR page
GET  /patients/:id/assessment       → Assessment form
GET  /patients/:id/wounds           → Wound documentation
GET  /patients/:id/tasks            → Task checklist
GET  /patients/:id/handoff          → Shift handoff form
GET  /messages                      → Messaging page
GET  /reports                       → Reports page

# API (SSE + JSON responses)
# Vitals
POST /api/patients/:id/vitals               → Record vitals
GET  /api/patients/:id/vitals/stream         → SSE: live vitals updates

# MAR
POST /api/patients/:id/mar                   → Administer medication
GET  /api/patients/:id/mar/stream            → SSE: MAR updates

# Assessment
POST /api/patients/:id/assessment            → Save assessment
POST /api/patients/:id/pain                  → Record pain assessment

# I&O
POST /api/patients/:id/io                    → Record intake/output
GET  /api/patients/:id/io/stream             → SSE: I&O totals

# Wounds
POST /api/patients/:id/wounds                → Document wound
POST /api/patients/:id/wounds/:wid/photo     → Upload wound photo

# Tasks
GET  /api/tasks/due                          → SSE: due/overdue tasks
POST /api/tasks/:id/complete                 → Complete a task
POST /api/patients/:id/tasks                 → Create task

# Handoff
POST /api/patients/:id/handoff               → Save shift report
GET  /api/patients/:id/handoff/prefill       → Get pre-filled SBAR data

# Communication
POST /api/messages                           → Send message
GET  /api/messages/stream                    → SSE: live messages
POST /api/notifications                      → Log physician notification

# Patient Management
POST /api/patients                           → Admit patient
PUT  /api/patients/:id                       → Update patient info
POST /api/patients/:id/discharge             → Discharge patient
GET  /api/census/stream                      → SSE: live census updates

# Safety
GET  /api/medications/interactions           → Check drug interactions
GET  /api/patients/:id/audit                 → View audit trail

# Reporting
GET  /api/patients/:id/shift-summary         → Generate shift summary
GET  /api/patients/:id/export                → Export patient data (CSV/PDF)
GET  /api/reports/abnormals                  → Abnormal values report
```

---

## Environment Variables

```bash
# Server
PORT=8080
ENV=development                    # development | production

# Database
DB_DRIVER=sqlite                   # sqlite | postgres
DB_DSN=./data/chartattack.db       # SQLite path or Postgres connection string

# Storage
PHOTO_STORAGE_PATH=./data/photos

# Security
SESSION_SECRET=<random-secret>
AUDIT_ENABLED=true

# Optional: External drug interaction API
DRUG_API_URL=
DRUG_API_KEY=
```

---

## Development Commands

```bash
# Run the server
go run ./cmd/server

# Run migrations
go run migrate.go up

# Generate sqlc code
sqlc generate

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...
```
