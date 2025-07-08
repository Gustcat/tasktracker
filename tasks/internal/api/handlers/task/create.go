package task

import (
	"errors"
	"fmt"
	"github.com/Gustcat/task-server/internal/api/handlers/converter"
	"github.com/Gustcat/task-server/internal/api/handlers/dto"
	"github.com/Gustcat/task-server/internal/lib/response"
	"github.com/Gustcat/task-server/internal/logger"
	"github.com/Gustcat/task-server/internal/repository"
	"github.com/Gustcat/task-server/internal/service"
	"github.com/Gustcat/task-server/internal/validation"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

func (h *Handler) Create(c *gin.Context) {
	const op = "handlers.task.Create"

	ctx := c.Request.Context()
	log := logger.LogFromContextAddOP(ctx, op)

	var requestTask dto.CreateTaskRequest

	log.Debug("Receive request for create task")
	if err := c.ShouldBindJSON(&requestTask); err != nil {
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
	log.Debug("Parsed create successfully", slog.Any("task", requestTask))

	task := converter.DTOToTask(&requestTask)
	id, err := h.service.Create(ctx, task)
	if errors.Is(err, repository.ErrTaskExists) {
		log.Error("Get error", slog.String("error", err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(fmt.Sprintf(
			"Task with title %s already exists", task.Title)))
		return
	}
	if errors.Is(err, service.ErrUserNotAllowed) {
		log.Error("User is not allowed to create task", slog.String("error", err.Error()))
		c.AbortWithStatusJSON(http.StatusForbidden, response.Error(err.Error()))
		return
	}
	if err != nil {
		log.Error("Failed to create task", slog.String("error", err.Error()))
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error("Failed to create task"))
		return
	}

	log.Info("Person created", slog.Int64("id", id))

	createResp := &dto.IdResponse{ID: id}
	c.JSON(http.StatusCreated, response.OK(createResp))
}
