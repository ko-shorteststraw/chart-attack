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

type MedicationController struct {
	medicationService appinterfaces.MedicationService
	patientService    appinterfaces.PatientService
	templates         map[string]*template.Template
}

func NewMedicationController(ms appinterfaces.MedicationService, ps appinterfaces.PatientService, tmpl map[string]*template.Template) *MedicationController {
	return &MedicationController{
		medicationService: ms,
		patientService:    ps,
		templates:         tmpl,
	}
}

func (mc *MedicationController) RegisterRoutes(e *echo.Echo) {
	e.GET("/patients/:id/mar", mc.MARPage)
	e.POST("/api/patients/:id/mar", mc.AdministerMedication)
	e.GET("/api/patients/:id/mar/stream", mc.StreamMAR)
}

func (mc *MedicationController) MARPage(c echo.Context) error {
	patientId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid patient ID")
	}

	patientResult, err := mc.patientService.GetPatientById(c.Request().Context(), query.GetPatientByIdQuery{PatientId: patientId})
	if err != nil {
		slog.Error("failed to get patient", "error", err)
		return c.String(http.StatusNotFound, "Patient not found")
	}

	marResult, err := mc.medicationService.GetPatientMAR(c.Request().Context(), query.GetPatientMARQuery{PatientId: patientId})
	if err != nil {
		slog.Error("failed to get MAR", "error", err)
	}

	medsResult, err := mc.medicationService.GetAllMedications(c.Request().Context(), query.GetAllMedicationsQuery{})
	if err != nil {
		slog.Error("failed to get medications", "error", err)
	}

	data := map[string]any{
		"Patient": patientResult.Patient,
	}
	if marResult != nil {
		data["Entries"] = marResult.Entries
	}
	if medsResult != nil {
		data["Medications"] = medsResult.Medications
	}

	return mc.templates["mar.html"].ExecuteTemplate(c.Response().Writer, "base", data)
}

func (mc *MedicationController) AdministerMedication(c echo.Context) error {
	patientId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid patient ID")
	}

	var signals struct {
		MedicationId string `json:"medicationId"`
		Dose         string `json:"dose"`
		Route        string `json:"route"`
	}
	if err := datastar.ReadSignals(c.Request(), &signals); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request")
	}

	medicationId, err := uuid.Parse(signals.MedicationId)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid medication ID")
	}

	demoUserId := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	cmd := command.AdministerMedicationCommand{
		IdempotencyKey: uuid.New().String(),
		PatientId:      patientId,
		MedicationId:   medicationId,
		ScheduledTime:  time.Now().UTC(),
		Dose:           signals.Dose,
		Route:          signals.Route,
		AdministeredBy: demoUserId,
	}

	if _, err := mc.medicationService.AdministerMedication(c.Request().Context(), cmd); err != nil {
		slog.Error("failed to administer medication", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to administer medication")
	}

	return mc.streamMARForPatient(c, patientId)
}

func (mc *MedicationController) StreamMAR(c echo.Context) error {
	patientId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid patient ID")
	}
	return mc.streamMARForPatient(c, patientId)
}

func (mc *MedicationController) streamMARForPatient(c echo.Context, patientId uuid.UUID) error {
	result, err := mc.medicationService.GetPatientMAR(c.Request().Context(), query.GetPatientMARQuery{PatientId: patientId})
	if err != nil {
		return err
	}
	sse := datastar.NewSSE(c.Response().Writer, c.Request())
	html, err := renderTemplate(mc.templates["mar.html"], "mar_table_body", map[string]any{
		"Entries": result.Entries,
	})
	if err != nil {
		return err
	}
	sse.PatchElements(html)
	return nil
}
