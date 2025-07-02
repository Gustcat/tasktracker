package task

import (
	"github.com/Gustcat/task-server/internal/repository"
	"github.com/Gustcat/task-server/internal/service"
)

type Serv struct {
	taskRepo   repository.TaskRepository
	authClient service.AuthService
}

func NewService(taskRepo repository.TaskRepository, authClient service.AuthService) service.TaskService {
	return &Serv{
		taskRepo:   taskRepo,
		authClient: authClient,
	}
}
