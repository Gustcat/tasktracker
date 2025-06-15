package task

import (
	"github.com/Gustcat/task-server/internal/service"
)

type Handler struct {
	service service.TaskService
}

func NewHandler(service service.TaskService) *Handler {
	return &Handler{service: service}
}
