package converter

import (
	"database/sql"
	"github.com/Gustcat/task-server/internal/model"
	modelRepo "github.com/Gustcat/task-server/internal/repository/model"
)

func TaskToRepo(task *model.TaskCreate) *modelRepo.TaskCreateDB {
	status := model.StatusNew
	if task.Status != nil {
		status = *task.Status
	}

	return &modelRepo.TaskCreateDB{
		Title:       task.Title,
		Description: pointerToSQL(task.Description),
		Status:      status,
		Author:      task.Author,
		Operator:    pointerToSQL(task.Operator),
		DueDate:     pointerToSQL(task.DueDate),
		CompletedAt: pointerToSQL(task.CompletedAt),
	}
}

func TaskUpdateToRepo(task *model.TaskUpdate) *modelRepo.TaskUpdateDB {
	return &modelRepo.TaskUpdateDB{
		Title:       pointerToSQL(task.Title),
		Description: pointerToSQL(task.Description),
		Status:      pointerToSQL(task.Status),
		Operator:    pointerToSQL(task.Operator),
		DueDate:     pointerToSQL(task.DueDate),
		CompletedAt: pointerToSQL(task.CompletedAt),
		UpdatedAt:   task.UpdatedAt,
	}
}

func pointerToSQL[T any](pointer *T) sql.Null[T] {
	if pointer == nil {
		return sql.Null[T]{Valid: false}
	}
	return sql.Null[T]{V: *pointer, Valid: true}
}

func RepoToTask(task *modelRepo.TaskDB) *model.Task {
	return &model.Task{
		ID:          task.ID,
		Title:       task.Title,
		Description: SQLToPointer(task.Description),
		Status:      task.Status,
		Author:      task.Author,
		Operator:    SQLToPointer(task.Operator),
		DueDate:     SQLToPointer(task.DueDate),
		CompletedAt: SQLToPointer(task.CompletedAt),
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   SQLToPointer(task.UpdatedAt),
	}
}

func RepoToFullTask(task *modelRepo.FullTaskDB) *model.FullTask {
	return &model.FullTask{
		Watchers: task.Watchers,
		Task:     *RepoToTask(&task.TaskDB),
	}
}

func SQLToPointer[T any](sqlField sql.Null[T]) *T {
	if sqlField.Valid {
		return &sqlField.V
	}
	return nil
}
