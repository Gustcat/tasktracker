package converter

import (
	"github.com/Gustcat/task-server/internal/api/handlers/dto"
	"github.com/Gustcat/task-server/internal/model"
	"time"
)

func DTOToTask(createTask *dto.CreateTaskRequest) *model.TaskCreate {
	return &model.TaskCreate{
		Title:       createTask.Title,
		Description: createTask.Description,
		Status:      createTask.Status,
		WatchSelf:   createTask.WatchSelf,
		Operator:    createTask.Operator,
		DueDate:     (*time.Time)(createTask.DueDate),
	}
}

func UpdateDTOToTaskUpdate(updateTask *dto.UpdateTaskRequest) *model.TaskUpdate {
	return &model.TaskUpdate{
		Title:       updateTask.Title,
		Description: updateTask.Description,
		Status:      updateTask.Status,
		Operator:    updateTask.Operator,
		WatchSelf:   updateTask.WatchSelf,
		DueDate:     (*time.Time)(updateTask.DueDate),
	}
}

func TaskToDTO(task *model.Task) *dto.TaskResponse {
	return &dto.TaskResponse{
		Title:           task.Title,
		Description:     task.Description,
		Status:          task.Status,
		Operator:        task.Operator,
		DueDate:         (*dto.Date)(task.DueDate),
		CreatedAt:       task.CreatedAt,
		UpdatedAt:       task.UpdatedAt,
		CompletedAt:     task.CompletedAt,
		Author:          task.Author,
		ID:              task.ID,
		AuthorDeleted:   task.AuthorDeleted,
		OperatorDeleted: task.OperatorDeleted,
	}
}

func FullTaskToDTO(task *model.FullTask) *dto.FullTaskResponse {
	return &dto.FullTaskResponse{
		TaskResponse: *TaskToDTO(&task.Task),
		Watchers:     task.Watchers,
	}
}
