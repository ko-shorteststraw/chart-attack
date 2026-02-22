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

type TaskController struct {
	taskService    appinterfaces.TaskService
	patientService appinterfaces.PatientService
	templates      map[string]*template.Template
}

func NewTaskController(ts appinterfaces.TaskService, ps appinterfaces.PatientService, tmpl map[string]*template.Template) *TaskController {
	return &TaskController{
		taskService:    ts,
		patientService: ps,
		templates:      tmpl,
	}
}

func (tc *TaskController) RegisterRoutes(e *echo.Echo) {
	e.GET("/patients/:id/tasks", tc.TasksPage)
	e.POST("/api/patients/:id/tasks", tc.CreateTask)
	e.POST("/api/tasks/:id/complete", tc.CompleteTask)
	e.GET("/api/patients/:id/tasks/stream", tc.StreamTasks)
}

func (tc *TaskController) TasksPage(c echo.Context) error {
	patientId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid patient ID")
	}

	patientResult, err := tc.patientService.GetPatientById(c.Request().Context(), query.GetPatientByIdQuery{PatientId: patientId})
	if err != nil {
		slog.Error("failed to get patient", "error", err)
		return c.String(http.StatusNotFound, "Patient not found")
	}

	tasksResult, err := tc.taskService.GetPatientTasks(c.Request().Context(), query.GetPatientTasksQuery{PatientId: patientId})
	if err != nil {
		slog.Error("failed to get tasks", "error", err)
	}

	data := map[string]any{
		"Patient": patientResult.Patient,
	}
	if tasksResult != nil {
		data["Tasks"] = tasksResult.Tasks
	}

	return tc.templates["tasks.html"].ExecuteTemplate(c.Response().Writer, "base", data)
}

func (tc *TaskController) CreateTask(c echo.Context) error {
	patientId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid patient ID")
	}

	var signals struct {
		TaskTitle    string `json:"taskTitle"`
		TaskCategory string `json:"taskCategory"`
		TaskPriority string `json:"taskPriority"`
		TaskDueAt    string `json:"taskDueAt"`
		TaskNotes    string `json:"taskNotes"`
	}
	if err := datastar.ReadSignals(c.Request(), &signals); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request")
	}

	dueAt := time.Now().Add(2 * time.Hour).UTC()
	if signals.TaskDueAt != "" {
		if d, err := time.Parse("2006-01-02T15:04", signals.TaskDueAt); err == nil {
			dueAt = d
		}
	}

	demoUserId := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	cmd := command.CreateTaskCommand{
		IdempotencyKey: uuid.New().String(),
		PatientId:      patientId,
		AssignedTo:     demoUserId,
		Title:          signals.TaskTitle,
		Category:       signals.TaskCategory,
		DueAt:          dueAt,
		Priority:       signals.TaskPriority,
		Notes:          signals.TaskNotes,
		CreatedBy:      demoUserId,
	}

	if _, err := tc.taskService.CreateTask(c.Request().Context(), cmd); err != nil {
		slog.Error("failed to create task", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to create task")
	}

	return tc.streamTasksForPatient(c, patientId)
}

func (tc *TaskController) CompleteTask(c echo.Context) error {
	taskId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid task ID")
	}

	demoUserId := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	cmd := command.CompleteTaskCommand{TaskId: taskId, CompletedBy: demoUserId}
	if err := tc.taskService.CompleteTask(c.Request().Context(), cmd); err != nil {
		slog.Error("failed to complete task", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to complete task")
	}

	// We need the patient ID to stream back — get the task to find it
	tasksResult, err := tc.taskService.GetPatientTasks(c.Request().Context(), query.GetPatientTasksQuery{PatientId: uuid.Nil})
	if err != nil {
		// Fallback: just send a simple success
		sse := datastar.NewSSE(c.Response().Writer, c.Request())
		sse.PatchElements(`<div id="tasks-items"><p>Task completed. Refresh to see updated list.</p></div>`)
		return nil
	}
	_ = tasksResult

	sse := datastar.NewSSE(c.Response().Writer, c.Request())
	sse.PatchElements(`<div id="tasks-items"><p class="text-muted">Task completed. Refresh to see updated list.</p></div>`)
	return nil
}

func (tc *TaskController) StreamTasks(c echo.Context) error {
	patientId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid patient ID")
	}
	return tc.streamTasksForPatient(c, patientId)
}

func (tc *TaskController) streamTasksForPatient(c echo.Context, patientId uuid.UUID) error {
	result, err := tc.taskService.GetPatientTasks(c.Request().Context(), query.GetPatientTasksQuery{PatientId: patientId})
	if err != nil {
		return err
	}
	sse := datastar.NewSSE(c.Response().Writer, c.Request())
	html, err := renderTemplate(tc.templates["tasks.html"], "tasks_list_body", map[string]any{
		"Tasks": result.Tasks,
	})
	if err != nil {
		return err
	}
	sse.PatchElements(html)
	return nil
}
