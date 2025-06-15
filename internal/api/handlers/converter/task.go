package converter

import (
	"github.com/Gustcat/task-server/internal/api/handlers/dto"
	"github.com/Gustcat/task-server/internal/model"
)

func DTOToTask(createTask *dto.CreateTaskRequest) *model.Task {
	return &model.Task{
		Title:       createTask.Title,
		Description: createTask.Description,
		Status:      createTask.Status,
		Operator:    createTask.Operator,
		DueDate:     createTask.DueDate,
	}
}
