package mapper

import (
	"github.com/kendall/chart-attack/internal/application/common"
	"github.com/kendall/chart-attack/internal/domain/entities"
)

func ToTaskResult(t *entities.Task) common.TaskResult {
	return common.TaskResult{
		Id:          t.Id.String(),
		PatientId:   t.PatientId.String(),
		AssignedTo:  t.AssignedTo.String(),
		Title:       t.Title,
		Category:    t.Category,
		DueAt:       t.DueAt,
		CompletedAt: t.CompletedAt,
		Priority:    t.Priority,
		Recurring:   t.Recurring,
		Notes:       t.Notes,
		IsCompleted: t.IsCompleted(),
	}
}

func ToTaskResults(tasks []*entities.Task) []common.TaskResult {
	results := make([]common.TaskResult, 0, len(tasks))
	for _, t := range tasks {
		results = append(results, ToTaskResult(t))
	}
	return results
}
