package converter

import (
	"github.com/Gustcat/task-server/internal/api/handlers/dto"
	"github.com/Gustcat/task-server/internal/model"
	"time"
)

func DTOToTask(createTask *dto.CreateTaskRequest) *model.Task {
	var t *time.Time
	if createTask.DueDate != nil {
		tt := time.Time(*createTask.DueDate)
		t = &tt
	}

	return &model.Task{
		Title:       createTask.Title,
		Description: createTask.Description,
		Status:      createTask.Status,
		Operator:    createTask.Operator,
		DueDate:     t,
	}
}

func TaskToDTO(task *model.Task) *dto.TaskResponse {
	return &dto.TaskResponse{
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		Operator:    task.Operator,
		DueDate:     task.DueDate,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
		CompletedAt: task.CompletedAt,
		Author:      task.Author,
		ID:          task.ID,
	}
}
