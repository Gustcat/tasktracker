package task

import (
	"errors"
	"fmt"
	"github.com/Gustcat/task-server/internal/api/handlers/dto"
	"github.com/Gustcat/task-server/internal/lib/response"
	"github.com/Gustcat/task-server/internal/logger"
	"github.com/Gustcat/task-server/internal/repository"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

func (h *Handler) Delete(c *gin.Context) {
	const op = "handlers.task.Delete"

	ctx := c.Request.Context()
	log := logger.LogFromContextAddOP(ctx, op)

	var idUri dto.IdUri

	err := c.ShouldBindUri(&idUri)
	if err != nil {
		log.Error("Invalid id parameter", slog.Int64("id", idUri.ID))
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error("invalid id url-parameter"))
		return
	}

	err = h.service.Delete(ctx, idUri.ID)
	if errors.Is(err, repository.ErrTaskNotFound) {
		log.Error(err.Error(), slog.Int64("id", idUri.ID))
		c.AbortWithStatusJSON(http.StatusNotFound, response.Error(fmt.Sprintf("%v, id - %d", err, idUri.ID)))
		return
	}

	if err != nil {
		log.Error("Failed to delete task")
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error(fmt.Sprintf("%v, id - %d", err, idUri.ID)))
		return
	}

	c.Status(http.StatusNoContent)

}
