package task

import (
	"context"
	"errors"
	"fmt"
	"github.com/Gustcat/task-server/internal/converter"
	"github.com/Gustcat/task-server/internal/lib/ctxutils"
	"github.com/Gustcat/task-server/internal/model"
	modelRepo "github.com/Gustcat/task-server/internal/repository/model"
	"github.com/Gustcat/task-server/internal/service"
	"time"
)

func (s *Serv) Update(ctx context.Context, id int64, task *model.TaskUpdate) (*model.Task, error) {
	const op = "service.task.Update"

	currentUser, err := ctxutils.UserFromContext(ctx)
	if currentUser == nil {
		return nil, fmt.Errorf("%w: %s", err, op)
	}

	if isEmptyTaskUpdate(task) && task.WatchSelf == nil {
		emptyUpdateRequest := errors.New("empty update request")
		return nil, fmt.Errorf("%s: %w", op, emptyUpdateRequest)
	}

	if currentUser.Role == model.USER {

		if task.Operator != nil && currentUser.ID != *task.Operator {
			return nil, fmt.Errorf("%w to assign anyone other than himself", service.ErrUserNotAllowed)
		}
	}

	task.UpdatedAt = time.Now()

	if task.Status != nil {
		if *task.Status == model.StatusDone {
			task.CompletedAt = &task.UpdatedAt
		} else {
			task.CompletedAt = nil
		}
	}

	var taskRepo *modelRepo.TaskDB
	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		if !isEmptyTaskUpdate(task) {
			taskRepo, errTx = s.taskRepo.Update(ctx, id, converter.TaskUpdateToRepo(task))
			if errTx != nil {
				return errTx
			}
		}

		if task.WatchSelf != nil {
			if *task.WatchSelf {
				errTx = s.watcherRepo.Add(ctx, id, currentUser.ID)
				if errTx != nil {
					return errTx
				}
				return nil
			} else {
				errTx = s.watcherRepo.Remove(ctx, id, currentUser.ID)
				if errTx != nil {
					return errTx
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return converter.RepoToTask(taskRepo), nil
}

func isEmptyTaskUpdate(t *model.TaskUpdate) bool {
	return t.Title == nil &&
		t.Description == nil &&
		t.Status == nil &&
		t.Operator == nil &&
		t.DueDate == nil
}
