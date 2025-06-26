package task

import (
	"github.com/Gustcat/task-server/internal/repository"
	"github.com/Gustcat/task-server/internal/service"
)

type Serv struct {
	taskRepo repository.TaskRepository
}

func NewService(taskRepo repository.TaskRepository) service.TaskService {
	return &Serv{
		taskRepo: taskRepo,
	}
}
