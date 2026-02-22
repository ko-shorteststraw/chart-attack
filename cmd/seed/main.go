package main

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/kendall/chart-attack/internal/domain/entities"
	"github.com/kendall/chart-attack/internal/infrastructure/db/sqlite"
)

func main() {
	dbDSN := os.Getenv("DB_DSN")
	if dbDSN == "" {
		dbDSN = "./data/chartattack.db"
	}

	db, err := sqlite.Open(dbDSN)
	if err != nil {
		slog.Error("failed to open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	migrationsDir := filepath.Join("migrations", "sqlite")
	if err := sqlite.RunMigrations(db, migrationsDir); err != nil {
		slog.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()
	userRepo := sqlite.NewSqlcUserRepository(db)
	patientRepo := sqlite.NewSqlcPatientRepository(db)

	// Create demo nurse
	nurse, err := entities.NewUser("jsmith", "Jane Smith, RN", "RN", "MedSurg", "BADGE001")
	if err != nil {
		slog.Error("failed to create nurse", "error", err)
		os.Exit(1)
	}
	validatedNurse, _ := entities.NewValidatedUser(nurse)
	if err := userRepo.Save(ctx, validatedNurse); err != nil {
		slog.Warn("nurse may already exist", "error", err)
	} else {
		slog.Info("created demo nurse", "id", nurse.Id, "name", nurse.FullName)
	}

	// Create sample patients
	patients := []struct {
		mrn, first, last, room, code, iso string
		allergies, diagnoses              []string
		fallRisk                          bool
		dob                               time.Time
	}{
		{
			mrn: "MRN001", first: "John", last: "Doe", room: "401-A",
			code: "Full Code", iso: "None", fallRisk: false,
			allergies: []string{"Penicillin"}, diagnoses: []string{"CHF", "HTN"},
			dob: time.Date(1955, 3, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			mrn: "MRN002", first: "Mary", last: "Johnson", room: "402-B",
			code: "DNR", iso: "Contact", fallRisk: true,
			allergies: []string{"Sulfa", "Latex"}, diagnoses: []string{"COPD", "DM2"},
			dob: time.Date(1948, 7, 22, 0, 0, 0, 0, time.UTC),
		},
		{
			mrn: "MRN003", first: "Robert", last: "Williams", room: "403-A",
			code: "Full Code", iso: "Droplet", fallRisk: false,
			allergies: []string{}, diagnoses: []string{"Pneumonia"},
			dob: time.Date(1970, 11, 8, 0, 0, 0, 0, time.UTC),
		},
		{
			mrn: "MRN004", first: "Patricia", last: "Brown", room: "404-A",
			code: "Comfort Care", iso: "None", fallRisk: true,
			allergies: []string{"Morphine", "Codeine"}, diagnoses: []string{"Pancreatic CA", "Pain"},
			dob: time.Date(1940, 1, 30, 0, 0, 0, 0, time.UTC),
		},
		{
			mrn: "MRN005", first: "James", last: "Davis", room: "405-B",
			code: "Full Code", iso: "None", fallRisk: false,
			allergies: []string{}, diagnoses: []string{"Post-op appendectomy"},
			dob: time.Date(1990, 5, 12, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, p := range patients {
		patient, err := entities.NewPatient(p.mrn, p.first, p.last, p.dob, p.room, p.allergies, p.diagnoses, p.code, p.iso, p.fallRisk)
		if err != nil {
			slog.Error("failed to create patient", "mrn", p.mrn, "error", err)
			continue
		}
		patient.AssignNurse(nurse.Id)
		validated, _ := entities.NewValidatedPatient(patient)
		if err := patientRepo.Save(ctx, validated); err != nil {
			slog.Warn("patient may already exist", "mrn", p.mrn, "error", err)
		} else {
			slog.Info("created patient", "mrn", p.mrn, "name", patient.FullName(), "room", p.room)
		}
	}

	slog.Info("seed complete")
}
