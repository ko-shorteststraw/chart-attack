package main

import (
	"html/template"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/kendall/chart-attack/internal/application/services"
	"github.com/kendall/chart-attack/internal/infrastructure/db/sqlite"
	"github.com/kendall/chart-attack/internal/interface/api/rest"
)

func main() {
	// Config
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	dbDSN := os.Getenv("DB_DSN")
	if dbDSN == "" {
		dbDSN = "./data/chartattack.db"
	}

	// Database
	db, err := sqlite.Open(dbDSN)
	if err != nil {
		slog.Error("failed to open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Run migrations
	migrationsDir := filepath.Join("migrations", "sqlite")
	if err := sqlite.RunMigrations(db, migrationsDir); err != nil {
		slog.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}
	slog.Info("database migrations complete")

	// Repositories
	userRepo := sqlite.NewSqlcUserRepository(db)
	_ = userRepo
	patientRepo := sqlite.NewSqlcPatientRepository(db)
	vitalSignRepo := sqlite.NewSqlcVitalSignRepository(db)
	auditRepo := sqlite.NewSqlcAuditRepository(db)
	idempotencyRepo := sqlite.NewSqlcIdempotencyRepository(db)
	medicationRepo := sqlite.NewSqlcMedicationRepository(db)
	marRepo := sqlite.NewSqlcMARRepository(db)
	vaccinationRepo := sqlite.NewSqlcVaccinationRepository(db)
	healthHistoryRepo := sqlite.NewSqlcHealthHistoryRepository(db)
	emergencyContactRepo := sqlite.NewSqlcEmergencyContactRepository(db)
	taskRepo := sqlite.NewSqlcTaskRepository(db)

	// Services
	patientService := services.NewPatientService(patientRepo, emergencyContactRepo, auditRepo, idempotencyRepo)
	chartingService := services.NewChartingService(vitalSignRepo, auditRepo, idempotencyRepo)
	medicationService := services.NewMedicationService(medicationRepo, marRepo, auditRepo, idempotencyRepo)
	vaccinationService := services.NewVaccinationService(vaccinationRepo, auditRepo, idempotencyRepo)
	healthHistoryService := services.NewHealthHistoryService(healthHistoryRepo, auditRepo, idempotencyRepo)
	taskService := services.NewTaskService(taskRepo, auditRepo, idempotencyRepo)

	// Templates - each page gets its own template set (base + partials + page)
	funcMap := template.FuncMap{
		"deref": func(p *int) int {
			if p == nil {
				return 0
			}
			return *p
		},
		"deref64": func(p *float64) float64 {
			if p == nil {
				return 0
			}
			return *p
		},
	}

	sharedFiles := []string{
		"templates/layouts/base.html",
		"templates/partials/patient_card.html",
		"templates/partials/vital_sign_row.html",
		"templates/partials/diagnosis_list.html",
		"templates/partials/emergency_contact_row.html",
		"templates/partials/vaccination_row.html",
		"templates/partials/health_history_row.html",
		"templates/partials/mar_entry_row.html",
		"templates/partials/task_item.html",
	}

	templates := make(map[string]*template.Template)
	pageFiles := map[string]string{
		"dashboard.html": "templates/pages/dashboard.html",
		"admit.html":     "templates/pages/admit.html",
		"profile.html":   "templates/pages/profile.html",
		"vitals.html":    "templates/pages/vitals.html",
		"mar.html":       "templates/pages/mar.html",
		"tasks.html":     "templates/pages/tasks.html",
		"messages.html":  "templates/pages/messages.html",
		"reports.html":   "templates/pages/reports.html",
	}
	for name, pageFile := range pageFiles {
		files := append([]string{pageFile}, sharedFiles...)
		t, err := template.New(filepath.Base(pageFile)).Funcs(funcMap).ParseFiles(files...)
		if err != nil {
			slog.Error("failed to parse template", "name", name, "error", err)
			os.Exit(1)
		}
		templates[name] = t
	}

	// Echo
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Static files
	e.Static("/static", "static")

	// Controllers
	patientController := rest.NewPatientController(patientService, templates)
	patientController.RegisterRoutes(e)

	chartingController := rest.NewChartingController(chartingService, patientService, vaccinationService, healthHistoryService, templates)
	chartingController.RegisterRoutes(e)

	medicationController := rest.NewMedicationController(medicationService, patientService, templates)
	medicationController.RegisterRoutes(e)

	taskController := rest.NewTaskController(taskService, patientService, templates)
	taskController.RegisterRoutes(e)

	// Start
	slog.Info("starting server", "port", port)
	if err := e.Start(":" + port); err != nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}
