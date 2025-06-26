package service

import (
	"context"
	"github.com/Gustcat/task-server/internal/model"
)

type TaskService interface {
	Create(ctx context.Context, task *model.TaskCreate, author int64) (int64, error)
	Get(ctx context.Context, id int64) (*model.Task, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, id int64, task *model.TaskUpdate) (*model.Task, error)
	List(ctx context.Context) ([]*model.Task, error)
}
