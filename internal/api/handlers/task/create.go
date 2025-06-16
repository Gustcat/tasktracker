package task

import (
	"context"
	"errors"
	"fmt"
	"github.com/Gustcat/task-server/internal/api/handlers/converter"
	"github.com/Gustcat/task-server/internal/api/handlers/dto"
	"github.com/Gustcat/task-server/internal/lib/response"
	"github.com/Gustcat/task-server/internal/repository"
	"github.com/Gustcat/task-server/internal/validation"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

func (h *Handler) Create(ctx context.Context, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.task.Create"

		log := log.With(slog.String("op", op))

		var author int64 = 1 // TODO: получение автора из токена. Сделать общую функцию для получения юзера и роли для всех хэндлеров

		var requestTask dto.CreateTaskRequest

		log.Debug("Receive request for create task")
		if err := c.ShouldBindJSON(&requestTask); err != nil {
			if validateErrs, ok := err.(validator.ValidationErrors); ok {
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
		id, err := h.service.Create(ctx, task, author)
		if errors.Is(err, repository.ErrTaskExists) {
			log.Error("Get error", slog.String("error", err.Error()))
			c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(fmt.Sprintf(
				"Task with title %s %s already exists", task.Title)))
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
}
