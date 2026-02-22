package rest

import (
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/starfederation/datastar-go/datastar"

	"github.com/kendall/chart-attack/internal/application/command"
	appinterfaces "github.com/kendall/chart-attack/internal/application/interfaces"
	"github.com/kendall/chart-attack/internal/application/query"
)

type PatientController struct {
	patientService appinterfaces.PatientService
	templates      map[string]*template.Template
}

func NewPatientController(ps appinterfaces.PatientService, tmpl map[string]*template.Template) *PatientController {
	return &PatientController{
		patientService: ps,
		templates:      tmpl,
	}
}

func (pc *PatientController) RegisterRoutes(e *echo.Echo) {
	e.GET("/", pc.Dashboard)
	e.GET("/patients/new", pc.AdmitForm)
	e.GET("/messages", pc.MessagesPage)
	e.GET("/reports", pc.ReportsPage)
	e.POST("/api/patients", pc.AdmitPatient)
	e.GET("/api/census/stream", pc.StreamCensus)
}

func (pc *PatientController) Dashboard(c echo.Context) error {
	result, err := pc.patientService.GetPatientCensus(c.Request().Context(), query.GetPatientCensusQuery{})
	if err != nil {
		slog.Error("failed to get census", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to load census")
	}

	data := map[string]any{
		"Patients": result.Patients,
	}
	return pc.templates["dashboard.html"].ExecuteTemplate(c.Response().Writer, "base", data)
}

func (pc *PatientController) AdmitForm(c echo.Context) error {
	return pc.templates["admit.html"].ExecuteTemplate(c.Response().Writer, "base", nil)
}

func (pc *PatientController) MessagesPage(c echo.Context) error {
	return pc.templates["messages.html"].ExecuteTemplate(c.Response().Writer, "base", nil)
}

func (pc *PatientController) ReportsPage(c echo.Context) error {
	return pc.templates["reports.html"].ExecuteTemplate(c.Response().Writer, "base", nil)
}

func (pc *PatientController) AdmitPatient(c echo.Context) error {
	var signals struct {
		Mrn           string `json:"mrn"`
		FirstName     string `json:"firstName"`
		LastName      string `json:"lastName"`
		DateOfBirth   string `json:"dateOfBirth"`
		RoomBed       string `json:"roomBed"`
		CodeStatus    string `json:"codeStatus"`
		IsolationType string `json:"isolationType"`
	}

	if err := datastar.ReadSignals(c.Request(), &signals); err != nil {
		slog.Error("failed to read signals", "error", err)
		return c.String(http.StatusBadRequest, "Invalid request")
	}

	dob, err := time.Parse("2006-01-02", signals.DateOfBirth)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid date of birth")
	}

	demoUserId := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	cmd := command.AdmitPatientCommand{
		IdempotencyKey: uuid.New().String(),
		MRN:            signals.Mrn,
		FirstName:      signals.FirstName,
		LastName:        signals.LastName,
		DateOfBirth:    dob,
		RoomBed:        signals.RoomBed,
		CodeStatus:     signals.CodeStatus,
		IsolationType:  signals.IsolationType,
		AdmittedBy:     demoUserId,
	}

	result, err := pc.patientService.AdmitPatient(c.Request().Context(), cmd)
	if err != nil {
		slog.Error("failed to admit patient", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to admit patient")
	}

	sse := datastar.NewSSE(c.Response().Writer, c.Request())
	sse.Redirect("/patients/" + result.PatientId)
	return nil
}

func (pc *PatientController) StreamCensus(c echo.Context) error {
	sse := datastar.NewSSE(c.Response().Writer, c.Request())

	result, err := pc.patientService.GetPatientCensus(c.Request().Context(), query.GetPatientCensusQuery{})
	if err != nil {
		slog.Error("failed to get census", "error", err)
		return err
	}

	for _, p := range result.Patients {
		html, err := renderTemplate(pc.templates["dashboard.html"], "patient_card", p)
		if err != nil {
			slog.Error("failed to render patient card", "error", err)
			continue
		}
		sse.PatchElements(html)
	}

	return nil
}
