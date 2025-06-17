package task

import (
	"context"
	"github.com/Gustcat/task-server/internal/converter"
	"github.com/Gustcat/task-server/internal/model"
)

func (s *Serv) List(ctx context.Context) ([]*model.Task, error) {
	repoTasks, err := s.taskRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	tasks := make([]*model.Task, len(repoTasks))

	for i, repoTask := range repoTasks {
		task := converter.RepoToTask(repoTask)
		tasks[i] = task
	}

	return tasks, nil
}
