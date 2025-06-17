package service

import (
	"context"
	"github.com/Gustcat/task-server/internal/model"
)

type TaskService interface {
	Create(ctx context.Context, task *model.Task, author int64) (int64, error)
	Get(ctx context.Context, id int64) (*model.Task, error)
	Delete(ctx context.Context, id int64) error
}
