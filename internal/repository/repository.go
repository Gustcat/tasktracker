package repository

import (
	"context"
	"errors"
	modelrepo "github.com/Gustcat/task-server/internal/repository/model"
)

var (
	ErrTaskNotFound = errors.New("task not found")
	ErrTaskExists   = errors.New("task already exists")
)

type TaskRepository interface {
	Create(ctx context.Context, task *modelrepo.TaskCreateDB) (int64, error)
}
