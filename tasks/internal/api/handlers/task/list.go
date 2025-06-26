package task

import (
	"github.com/Gustcat/task-server/internal/api/handlers/converter"
	"github.com/Gustcat/task-server/internal/api/handlers/dto"
	"github.com/Gustcat/task-server/internal/lib/response"
	"github.com/Gustcat/task-server/internal/logger"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

func (h *Handler) List(c *gin.Context) {
	const op = "handlers.task.Delete"

	ctx := c.Request.Context()
	log := logger.LogFromContextAddOP(ctx, op)

	tasks, err := h.service.List(ctx)
	if err != nil {
		log.Error("Failed to list tasks", slog.String("error", err.Error()))
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error("Failed to list task"))
		return
	}

	tasksDto := make([]*dto.TaskResponse, len(tasks))

	for i, task := range tasks {
		taskDto := converter.TaskToDTO(task)
		tasksDto[i] = taskDto
	}

	c.JSON(http.StatusOK, response.OK(&tasksDto))
}
