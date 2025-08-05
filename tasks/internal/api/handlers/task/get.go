package task

import (
	"errors"
	"fmt"
	"github.com/Gustcat/task-server/internal/api/handlers/converter"
	"github.com/Gustcat/task-server/internal/api/handlers/dto"
	"github.com/Gustcat/task-server/internal/lib/response"
	"github.com/Gustcat/task-server/internal/logger"
	"github.com/Gustcat/task-server/internal/repository"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

func (h *Handler) Get(c *gin.Context) {
	const op = "handlers.task.Get"

	ctx := c.Request.Context()
	log := logger.LogFromContextAddOP(ctx, op)

	var idUri dto.IdUri

	err := c.ShouldBindUri(&idUri)
	if err != nil {
		log.Error("invalid url-parameter: id")
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error("invalid url-parameter: id"))
		return
	}

	task, err := h.service.Get(ctx, idUri.ID)
	if errors.Is(err, repository.ErrTaskNotFound) {
		log.Error(err.Error(), slog.Int64("id", idUri.ID))
		c.AbortWithStatusJSON(http.StatusNotFound, response.Error(fmt.Sprintf("%v, id - %d", err, idUri.ID)))
		return
	}

	if err != nil {
		log.Error("Failed to get task", slog.String("error", err.Error()))
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error("Failed to get task"))
		return
	}

	c.JSON(http.StatusOK, response.OK(converter.FullTaskToDTO(task)))
}
