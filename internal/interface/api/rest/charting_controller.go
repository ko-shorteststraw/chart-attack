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

type ChartingController struct {
	chartingService      appinterfaces.ChartingService
	patientService       appinterfaces.PatientService
	vaccinationService   appinterfaces.VaccinationService
	healthHistoryService appinterfaces.HealthHistoryService
	templates            map[string]*template.Template
}

func NewChartingController(
	cs appinterfaces.ChartingService,
	ps appinterfaces.PatientService,
	vs appinterfaces.VaccinationService,
	hhs appinterfaces.HealthHistoryService,
	tmpl map[string]*template.Template,
) *ChartingController {
	return &ChartingController{
		chartingService:      cs,
		patientService:       ps,
		vaccinationService:   vs,
		healthHistoryService: hhs,
		templates:            tmpl,
	}
}

func (cc *ChartingController) RegisterRoutes(e *echo.Echo) {
	// Profile page
	e.GET("/patients/:id", cc.ProfilePage)

	// Vitals
	e.GET("/patients/:id/vitals", cc.VitalsPage)
	e.POST("/api/patients/:id/vitals", cc.RecordVitals)
	e.GET("/api/patients/:id/vitals/stream", cc.StreamVitals)

	// Diagnoses
	e.POST("/api/patients/:id/diagnoses", cc.AddDiagnosis)
	e.POST("/api/patients/:id/diagnoses/remove", cc.RemoveDiagnosis)
	e.GET("/api/patients/:id/diagnoses/stream", cc.StreamDiagnoses)

	// Emergency Contacts
	e.POST("/api/patients/:id/emergency-contacts", cc.AddEmergencyContact)
	e.POST("/api/emergency-contacts/:id/remove", cc.RemoveEmergencyContact)
	e.GET("/api/patients/:id/emergency-contacts/stream", cc.StreamEmergencyContacts)

	// Vaccinations
	e.POST("/api/patients/:id/vaccinations", cc.AddVaccination)
	e.GET("/api/patients/:id/vaccinations/stream", cc.StreamVaccinations)

	// Health History
	e.POST("/api/patients/:id/health-history", cc.AddHealthHistory)
	e.GET("/api/patients/:id/health-history/stream", cc.StreamHealthHistory)
}

func (cc *ChartingController) ProfilePage(c echo.Context) error {
	patientId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid patient ID")
	}

	patientResult, err := cc.patientService.GetPatientById(c.Request().Context(), query.GetPatientByIdQuery{PatientId: patientId})
	if err != nil {
		slog.Error("failed to get patient", "error", err)
		return c.String(http.StatusNotFound, "Patient not found")
	}

	contactsResult, err := cc.patientService.GetPatientEmergencyContacts(c.Request().Context(), query.GetPatientEmergencyContactsQuery{PatientId: patientId})
	if err != nil {
		slog.Error("failed to get emergency contacts", "error", err)
	}

	vacResult, err := cc.vaccinationService.GetPatientVaccinations(c.Request().Context(), query.GetPatientVaccinationsQuery{PatientId: patientId})
	if err != nil {
		slog.Error("failed to get vaccinations", "error", err)
	}

	hhResult, err := cc.healthHistoryService.GetPatientHealthHistory(c.Request().Context(), query.GetPatientHealthHistoryQuery{PatientId: patientId})
	if err != nil {
		slog.Error("failed to get health history", "error", err)
	}

	data := map[string]any{
		"Patient": patientResult.Patient,
	}
	if contactsResult != nil {
		data["Contacts"] = contactsResult.Contacts
	}
	if vacResult != nil {
		data["Vaccinations"] = vacResult.Vaccinations
	}
	if hhResult != nil {
		data["HealthHistory"] = hhResult.Entries
	}

	return cc.templates["profile.html"].ExecuteTemplate(c.Response().Writer, "base", data)
}

func (cc *ChartingController) VitalsPage(c echo.Context) error {
	patientId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid patient ID")
	}

	patientResult, err := cc.patientService.GetPatientById(c.Request().Context(), query.GetPatientByIdQuery{PatientId: patientId})
	if err != nil {
		slog.Error("failed to get patient", "error", err)
		return c.String(http.StatusNotFound, "Patient not found")
	}

	vitalsResult, err := cc.chartingService.GetPatientVitals(c.Request().Context(), query.GetPatientVitalsQuery{PatientId: patientId})
	if err != nil {
		slog.Error("failed to get vitals", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to load vitals")
	}

	data := map[string]any{
		"Patient": patientResult.Patient,
		"Vitals":  vitalsResult.Vitals,
	}
	return cc.templates["vitals.html"].ExecuteTemplate(c.Response().Writer, "base", data)
}

func (cc *ChartingController) RecordVitals(c echo.Context) error {
	patientId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid patient ID")
	}

	var signals struct {
		Systolic  any    `json:"systolic"`
		Diastolic any    `json:"diastolic"`
		HR        any    `json:"hr"`
		Temp      any    `json:"temp"`
		TempRoute string `json:"tempRoute"`
		O2        any    `json:"o2"`
		Resp      any    `json:"resp"`
		Pain      any    `json:"pain"`
		Position  string `json:"position"`
		Notes     string `json:"notes"`
	}

	if err := datastar.ReadSignals(c.Request(), &signals); err != nil {
		slog.Error("failed to read signals", "error", err)
		return c.String(http.StatusBadRequest, "Invalid request")
	}

	demoUserId := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	cmd := command.RecordVitalsCommand{
		IdempotencyKey: uuid.New().String(),
		PatientId:      patientId,
		RecordedBy:     demoUserId,
		SystolicBP:     parseIntSignal(signals.Systolic),
		DiastolicBP:    parseIntSignal(signals.Diastolic),
		HeartRate:      parseIntSignal(signals.HR),
		Temperature:    parseFloatSignal(signals.Temp),
		TempRoute:      signals.TempRoute,
		OxygenSat:      parseIntSignal(signals.O2),
		Respirations:   parseIntSignal(signals.Resp),
		PainLevel:      parseIntSignal(signals.Pain),
		Position:       signals.Position,
		Notes:          signals.Notes,
	}

	_, err = cc.chartingService.RecordVitals(c.Request().Context(), cmd)
	if err != nil {
		slog.Error("failed to record vitals", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to record vitals")
	}

	sse := datastar.NewSSE(c.Response().Writer, c.Request())

	vitalsResult, err := cc.chartingService.GetPatientVitals(c.Request().Context(), query.GetPatientVitalsQuery{PatientId: patientId})
	if err != nil {
		slog.Error("failed to get vitals after record", "error", err)
		return err
	}

	html, err := renderTemplate(cc.templates["vitals.html"], "vitals_table_body", map[string]any{
		"Vitals": vitalsResult.Vitals,
	})
	if err != nil {
		slog.Error("failed to render vitals table", "error", err)
		return err
	}

	sse.PatchElements(html)
	return nil
}

func (cc *ChartingController) StreamVitals(c echo.Context) error {
	patientId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid patient ID")
	}

	sse := datastar.NewSSE(c.Response().Writer, c.Request())

	vitalsResult, err := cc.chartingService.GetPatientVitals(c.Request().Context(), query.GetPatientVitalsQuery{PatientId: patientId})
	if err != nil {
		slog.Error("failed to get vitals", "error", err)
		return err
	}

	html, err := renderTemplate(cc.templates["vitals.html"], "vitals_table_body", map[string]any{
		"Vitals": vitalsResult.Vitals,
	})
	if err != nil {
		slog.Error("failed to render vitals table", "error", err)
		return err
	}

	sse.PatchElements(html)
	return nil
}

// Diagnoses

func (cc *ChartingController) AddDiagnosis(c echo.Context) error {
	patientId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid patient ID")
	}

	var signals struct {
		NewDiagnosis string `json:"newDiagnosis"`
	}
	if err := datastar.ReadSignals(c.Request(), &signals); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request")
	}

	demoUserId := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	cmd := command.UpdateDiagnosesCommand{PatientId: patientId, Action: "add", Diagnosis: signals.NewDiagnosis, UpdatedBy: demoUserId}
	if err := cc.patientService.UpdateDiagnoses(c.Request().Context(), cmd); err != nil {
		slog.Error("failed to add diagnosis", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to add diagnosis")
	}

	return cc.streamDiagnosesForPatient(c, patientId)
}

func (cc *ChartingController) RemoveDiagnosis(c echo.Context) error {
	patientId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid patient ID")
	}

	var signals struct {
		RemoveDiagnosis string `json:"removeDiagnosis"`
	}
	if err := datastar.ReadSignals(c.Request(), &signals); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request")
	}

	demoUserId := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	cmd := command.UpdateDiagnosesCommand{PatientId: patientId, Action: "remove", Diagnosis: signals.RemoveDiagnosis, UpdatedBy: demoUserId}
	if err := cc.patientService.UpdateDiagnoses(c.Request().Context(), cmd); err != nil {
		slog.Error("failed to remove diagnosis", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to remove diagnosis")
	}

	return cc.streamDiagnosesForPatient(c, patientId)
}

func (cc *ChartingController) StreamDiagnoses(c echo.Context) error {
	patientId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid patient ID")
	}
	return cc.streamDiagnosesForPatient(c, patientId)
}

func (cc *ChartingController) streamDiagnosesForPatient(c echo.Context, patientId uuid.UUID) error {
	patientResult, err := cc.patientService.GetPatientById(c.Request().Context(), query.GetPatientByIdQuery{PatientId: patientId})
	if err != nil {
		return err
	}
	sse := datastar.NewSSE(c.Response().Writer, c.Request())
	html, err := renderTemplate(cc.templates["profile.html"], "diagnosis_list", map[string]any{
		"Patient": patientResult.Patient,
	})
	if err != nil {
		return err
	}
	sse.PatchElements(html)
	return nil
}

// Emergency Contacts

func (cc *ChartingController) AddEmergencyContact(c echo.Context) error {
	patientId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid patient ID")
	}

	var signals struct {
		EcName         string `json:"ecName"`
		EcRelationship string `json:"ecRelationship"`
		EcPhone        string `json:"ecPhone"`
		EcEmail        string `json:"ecEmail"`
	}
	if err := datastar.ReadSignals(c.Request(), &signals); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request")
	}

	demoUserId := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	cmd := command.AddEmergencyContactCommand{
		PatientId:    patientId,
		Name:         signals.EcName,
		Relationship: signals.EcRelationship,
		Phone:        signals.EcPhone,
		Email:        signals.EcEmail,
		AddedBy:      demoUserId,
	}

	if _, err := cc.patientService.AddEmergencyContact(c.Request().Context(), cmd); err != nil {
		slog.Error("failed to add emergency contact", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to add contact")
	}

	return cc.streamEmergencyContactsForPatient(c, patientId)
}

func (cc *ChartingController) RemoveEmergencyContact(c echo.Context) error {
	contactId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid contact ID")
	}

	demoUserId := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	cmd := command.RemoveEmergencyContactCommand{ContactId: contactId, RemovedBy: demoUserId}
	if err := cc.patientService.RemoveEmergencyContact(c.Request().Context(), cmd); err != nil {
		slog.Error("failed to remove emergency contact", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to remove contact")
	}

	// We can't easily know the patientId from just the contact ID, so redirect
	sse := datastar.NewSSE(c.Response().Writer, c.Request())
	sse.PatchElements(`<div id="ec-items"><p class="text-muted" style="margin:0;">Contact removed. Refresh to see updated list.</p></div>`)
	return nil
}

func (cc *ChartingController) StreamEmergencyContacts(c echo.Context) error {
	patientId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid patient ID")
	}
	return cc.streamEmergencyContactsForPatient(c, patientId)
}

func (cc *ChartingController) streamEmergencyContactsForPatient(c echo.Context, patientId uuid.UUID) error {
	result, err := cc.patientService.GetPatientEmergencyContacts(c.Request().Context(), query.GetPatientEmergencyContactsQuery{PatientId: patientId})
	if err != nil {
		return err
	}
	sse := datastar.NewSSE(c.Response().Writer, c.Request())
	html, err := renderTemplate(cc.templates["profile.html"], "emergency_contacts_list", map[string]any{
		"Contacts": result.Contacts,
	})
	if err != nil {
		return err
	}
	sse.PatchElements(html)
	return nil
}

// Vaccinations

func (cc *ChartingController) AddVaccination(c echo.Context) error {
	patientId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid patient ID")
	}

	var signals struct {
		VaccineName string `json:"vaccineName"`
		VaccineDate string `json:"vaccineDate"`
		VaccineLot  string `json:"vaccineLot"`
		VaccineSite string `json:"vaccineSite"`
	}
	if err := datastar.ReadSignals(c.Request(), &signals); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request")
	}

	dateAdministered := time.Now().UTC()
	if signals.VaccineDate != "" {
		if d, err := time.Parse("2006-01-02", signals.VaccineDate); err == nil {
			dateAdministered = d
		}
	}

	demoUserId := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	cmd := command.AddVaccinationCommand{
		IdempotencyKey:   uuid.New().String(),
		PatientId:        patientId,
		RecordedBy:       demoUserId,
		VaccineName:      signals.VaccineName,
		DateAdministered: dateAdministered,
		LotNumber:        signals.VaccineLot,
		Site:             signals.VaccineSite,
	}

	if _, err := cc.vaccinationService.AddVaccination(c.Request().Context(), cmd); err != nil {
		slog.Error("failed to add vaccination", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to add vaccination")
	}

	return cc.streamVaccinationsForPatient(c, patientId)
}

func (cc *ChartingController) StreamVaccinations(c echo.Context) error {
	patientId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid patient ID")
	}
	return cc.streamVaccinationsForPatient(c, patientId)
}

func (cc *ChartingController) streamVaccinationsForPatient(c echo.Context, patientId uuid.UUID) error {
	result, err := cc.vaccinationService.GetPatientVaccinations(c.Request().Context(), query.GetPatientVaccinationsQuery{PatientId: patientId})
	if err != nil {
		return err
	}
	sse := datastar.NewSSE(c.Response().Writer, c.Request())
	html, err := renderTemplate(cc.templates["profile.html"], "vaccinations_list", map[string]any{
		"Vaccinations": result.Vaccinations,
	})
	if err != nil {
		return err
	}
	sse.PatchElements(html)
	return nil
}

// Health History

func (cc *ChartingController) AddHealthHistory(c echo.Context) error {
	patientId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid patient ID")
	}

	var signals struct {
		HhCondition   string `json:"hhCondition"`
		HhStatus      string `json:"hhStatus"`
		HhDescription string `json:"hhDescription"`
	}
	if err := datastar.ReadSignals(c.Request(), &signals); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request")
	}

	demoUserId := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	cmd := command.AddHealthHistoryCommand{
		IdempotencyKey: uuid.New().String(),
		PatientId:      patientId,
		RecordedBy:     demoUserId,
		Condition:      signals.HhCondition,
		Status:         signals.HhStatus,
		Description:    signals.HhDescription,
	}

	if _, err := cc.healthHistoryService.AddHealthHistory(c.Request().Context(), cmd); err != nil {
		slog.Error("failed to add health history", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to add health history")
	}

	return cc.streamHealthHistoryForPatient(c, patientId)
}

func (cc *ChartingController) StreamHealthHistory(c echo.Context) error {
	patientId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid patient ID")
	}
	return cc.streamHealthHistoryForPatient(c, patientId)
}

func (cc *ChartingController) streamHealthHistoryForPatient(c echo.Context, patientId uuid.UUID) error {
	result, err := cc.healthHistoryService.GetPatientHealthHistory(c.Request().Context(), query.GetPatientHealthHistoryQuery{PatientId: patientId})
	if err != nil {
		return err
	}
	sse := datastar.NewSSE(c.Response().Writer, c.Request())
	html, err := renderTemplate(cc.templates["profile.html"], "health_history_list", map[string]any{
		"HealthHistory": result.Entries,
	})
	if err != nil {
		return err
	}
	sse.PatchElements(html)
	return nil
}
