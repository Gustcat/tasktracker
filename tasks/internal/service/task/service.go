package task

import (
	"github.com/Gustcat/task-server/internal/client/db"
	"github.com/Gustcat/task-server/internal/repository"
	"github.com/Gustcat/task-server/internal/service"
)

type Serv struct {
	taskRepo    repository.TaskRepository
	watcherRepo repository.WatcherRepository
	authClient  service.AuthService
	txManager   db.TxManager
}

func NewService(taskRepo repository.TaskRepository,
	watcherRepo repository.WatcherRepository,
	authClient service.AuthService,
	txManager db.TxManager) service.TaskService {
	return &Serv{
		taskRepo:    taskRepo,
		watcherRepo: watcherRepo,
		authClient:  authClient,
		txManager:   txManager,
	}
}
