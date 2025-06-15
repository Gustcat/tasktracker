package converter

import (
	"database/sql"
	"github.com/Gustcat/task-server/internal/model"
	modelRepo "github.com/Gustcat/task-server/internal/repository/model"
	"time"
)

func TaskToRepo(task *model.Task) *modelRepo.TaskCreateDB {
	return &modelRepo.TaskCreateDB{
		Title:       task.Title,
		Description: pointerToSQL[string](task.Description),
		Status:      task.Status,
		Author:      task.Author,
		Operator:    pointerToSQL[int64](task.Operator),
		DueDate:     pointerToSQL[time.Time](task.DueDate),
		CompletedAt: pointerToSQL[time.Time](task.CompletedAt),
	}
}

func pointerToSQL[T any](pointer *T) sql.Null[T] {
	if pointer == nil {
		return sql.Null[T]{Valid: false}
	}
	return sql.Null[T]{V: *pointer, Valid: true}
}
