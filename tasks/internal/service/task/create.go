package task

import (
	"context"
	"fmt"
	"github.com/Gustcat/task-server/internal/converter"
	"github.com/Gustcat/task-server/internal/model"
	"github.com/Gustcat/task-server/internal/service"
	"time"
)

func (s *Serv) Create(ctx context.Context, task *model.TaskCreate, authorId int64) (int64, error) {

	author, _, err := s.validateUsers(ctx, authorId, task.Operator)
	if err != nil {
		return 0, err
	}

	if author.Role == model.USER {
		if task.Operator != nil && authorId != *task.Operator {
			return 0, fmt.Errorf("%w to assign anyone other than himself", service.ErrUserNotAllowed)
		}
		if task.Status != nil && *task.Status == model.StatusDone {
			return 0, fmt.Errorf("%w to create task with `%s` status", service.ErrUserNotAllowed, model.StatusDone)
		}
	} // получить роль можно же из токена

	task.Author = authorId

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

	insertTask := converter.TaskToRepo(task)

	id, err := s.taskRepo.Create(ctx, insertTask)
	if err != nil {
		return 0, err
	}

	return id, nil
}
