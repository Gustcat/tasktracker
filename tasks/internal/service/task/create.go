package task

import (
	"context"
	"errors"
	"fmt"
	"github.com/Gustcat/task-server/internal/converter"
	"github.com/Gustcat/task-server/internal/lib/ctxutils"
	"github.com/Gustcat/task-server/internal/model"
	"github.com/Gustcat/task-server/internal/service"
	"time"
)

func (s *Serv) Create(ctx context.Context, task *model.TaskCreate) (int64, error) {
	author := ctxutils.UserFromContext(ctx)
	if author == nil {
		return 0, fmt.Errorf("%w: author", ErrUserNotFound) //TODO: обработать такую ошибку слоем выше
	}

	if author.Role == model.USER {
		if task.Operator != nil && author.ID != *task.Operator {
			return 0, fmt.Errorf("%w to assign anyone other than himself", service.ErrUserNotAllowed)
		}
		if task.Status != nil && *task.Status == model.StatusDone {
			return 0, fmt.Errorf("%w to create task with `%s` status", service.ErrUserNotAllowed, model.StatusDone)
		}
	}

	if task.Operator != nil {
		_, err := s.validateUser(ctx, *task.Operator)
		if errors.Is(err, ErrUserNotFound) {
			return 0, fmt.Errorf("%w: operator", ErrUserNotFound)
		}
		if err != nil {
			return 0, err
		}
	}

	if task.Status == nil {
		if task.Operator == nil {
			*task.Status = model.StatusNew
		} else {
			*task.Status = model.StatusTODO
		}
	} else {
		if *task.Status == model.StatusDone {
			*task.CompletedAt = time.Now()
		}
	}

	task.Author = author.ID
	insertTask := converter.TaskToRepo(task)

	id, err := s.taskRepo.Create(ctx, insertTask, task.WatchSelf)
	if err != nil {
		return 0, err
	}

	return id, nil
}
