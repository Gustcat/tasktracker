package task

import (
	"errors"
	"fmt"
	"github.com/Gustcat/task-server/internal/api/handlers/converter"
	"github.com/Gustcat/task-server/internal/api/handlers/dto"
	"github.com/Gustcat/task-server/internal/lib/response"
	"github.com/Gustcat/task-server/internal/logger"
	"github.com/Gustcat/task-server/internal/repository"
	"github.com/Gustcat/task-server/internal/validation"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

func (h *Handler) Update(c *gin.Context) {
	const op = "handlers.task.Update"

	ctx := c.Request.Context()
	log := logger.LogFromContextAddOP(ctx, op)

	var idUri dto.IdUri

	err := c.ShouldBindUri(&idUri)
	if err != nil {
		log.Error("Invalid id parameter", slog.Int64("id", idUri.ID))
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error("invalid id url-parameter"))
		return
	}

	var task dto.UpdateTaskRequest

	if err := c.ShouldBindJSON(&task); err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			errMsg := validation.ErrorMessage(validateErrs)
			log.Error("Validation failure", slog.String("error", errMsg))
			c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(errMsg))
			return
		}
		log.Error("Failed to parse request", slog.String("error", err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error("failed to parse request"))
		return
	}
	log.Debug("Parsed update successfully", slog.Any("task", task))

	updatedTask, err := h.service.Update(ctx, idUri.ID, converter.UpdateDTOToTaskUpdate(&task))
	if errors.Is(err, repository.ErrTaskNotFound) {
		log.Error(err.Error(), slog.Int64("id", idUri.ID))
		c.AbortWithStatusJSON(http.StatusNotFound, response.Error(fmt.Sprintf("%v, id - %d", err, idUri.ID)))
		return
	}

	if err != nil {
		log.Error("Failed to update task", slog.String("error", err.Error()))
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error("Failed to update task"))
		return
	}

	c.JSON(http.StatusOK, response.OK(converter.TaskToDTO(updatedTask)))
}
