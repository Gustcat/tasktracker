package task

import (
	"context"
	"errors"
	"fmt"
	"github.com/Gustcat/task-server/internal/converter"
	"github.com/Gustcat/task-server/internal/model"
	"time"
)

func (s *Serv) Update(ctx context.Context, id int64, task *model.TaskUpdate) (*model.Task, error) {
	const op = "service.task.Update"

	if isEmptyTaskUpdate(task) {
		emptyUpdateRequest := errors.New("empty update request")
		return nil, fmt.Errorf("%s: %w", op, emptyUpdateRequest)
	}

	task.UpdatedAt = time.Now()

	if task.Status != nil {
		if *task.Status == model.StatusDone {
			task.CompletedAt = &task.UpdatedAt
		} else {
			task.CompletedAt = nil
		}
	}

	taskRepo, err := s.taskRepo.Update(ctx, id, converter.TaskUpdateToRepo(task))
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
