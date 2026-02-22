package query

import (
	"github.com/google/uuid"
	"github.com/kendall/chart-attack/internal/application/common"
)

type GetPatientTasksQuery struct {
	PatientId uuid.UUID
}

type GetPatientTasksResult struct {
	Tasks []common.TaskResult
}
