package task

import (
	"github.com/Gustcat/task-server/internal/repository"
	"github.com/Gustcat/task-server/internal/service"
)

type Serv struct {
	taskRepo    repository.TaskRepository
	watcherRepo repository.WatcherRepository
	authClient  service.AuthService
}

func NewService(taskRepo repository.TaskRepository,
	watcherRepo repository.WatcherRepository,
	authClient service.AuthService) service.TaskService {
	return &Serv{
		taskRepo:    taskRepo,
		watcherRepo: watcherRepo,
		authClient:  authClient,
	}
}
