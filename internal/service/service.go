package service

import (
	"context"
	"github.com/Gustcat/task-server/internal/model"
)

type TaskService interface {
	Create(ctx context.Context, task *model.Task, author int64) (int64, error)
}
