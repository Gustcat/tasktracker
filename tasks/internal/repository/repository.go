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
	Get(ctx context.Context, id int64) (*modelrepo.TaskDB, error)
	GetWithWatchers(ctx context.Context, id int64) (*modelrepo.FullTaskDB, error)
	List(ctx context.Context) ([]*modelrepo.TaskDB, error)
	Update(ctx context.Context, id int64, task *modelrepo.TaskUpdateDB) (*modelrepo.TaskDB, error)
	Delete(ctx context.Context, id int64) error
	MarkAuthorDeleted(ctx context.Context, userID int64) ([]int64, error)
	MarkOperatorDeleted(ctx context.Context, userID int64) ([]int64, error)
}

type WatcherRepository interface {
	Add(ctx context.Context, taskID, userID int64) error
	Remove(ctx context.Context, taskID, userID int64) error
	DeleteUser(ctx context.Context, userID int64) error
}
