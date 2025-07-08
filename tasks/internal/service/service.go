package service

import (
	"context"
	"errors"
	"github.com/Gustcat/task-server/internal/model"
)

var (
	ErrUserNotAllowed = errors.New("USER role is not allowed")
)

type TaskService interface {
	Create(ctx context.Context, task *model.TaskCreate) (int64, error)
	Get(ctx context.Context, id int64) (*model.Task, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, id int64, task *model.TaskUpdate) (*model.Task, error)
	List(ctx context.Context) ([]*model.Task, error)
}

type AuthService interface {
	GetUser(ctx context.Context, id int64) (*model.User, error)
}
