package v1

import (
	"ai-assistants-catalog/internal/runs/app"
	"ai-assistants-catalog/internal/runs/domain"
)

func mapRun(run domain.Run) RunDTO {
	return RunDTO{
		ID:            run.ID,
		AssistantID:   run.AssistantID,
		AssistantName: run.AssistantName,
		CategoryID:    run.CategoryID,
		CategoryName:  run.CategoryName,
		UserID:        run.UserID,
		Model:         run.Model,
		UserPrompt:    run.UserPrompt,
		Output:        run.Output,
		Status:        run.Status,
		Error:         run.Error,
		CreatedAt:     run.CreatedAt,
	}
}

func mapRuns(runs []domain.Run) []RunDTO {
	result := make([]RunDTO, 0, len(runs))
	for _, run := range runs {
		result = append(result, mapRun(run))
	}

	return result
}

func mapListResult(result app.ListResult) RunsResponse {
	return RunsResponse{
		Runs: mapRuns(result.Runs),
		Pagination: PaginationDTO{
			Page:     result.Page,
			PageSize: result.PageSize,
			Total:    result.Total,
		},
	}
}
