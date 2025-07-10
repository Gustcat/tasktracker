package task

import (
	"context"
	"github.com/Gustcat/task-server/internal/converter"
	"github.com/Gustcat/task-server/internal/model"
)

func (s *Serv) Get(ctx context.Context, id int64) (*model.FullTask, error) {
	taskRepo, err := s.taskRepo.GetWithWatchers(ctx, id)
	if err != nil {
		return nil, err
	}

	return converter.RepoToFullTask(taskRepo), nil
}
