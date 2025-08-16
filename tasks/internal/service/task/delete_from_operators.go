package task

import (
	"context"
	"github.com/Gustcat/task-server/internal/logger"
	"log/slog"
)

func (s *Serv) DeleteUserFromOperators(ctx context.Context, userID int64) error {
	const op = "service.task.DeleteUserFromOperators"

	log := logger.LogFromContextAddOP(ctx, op)
	log.Info("DeleteUserFromOperators called")

	task_ids, err := s.taskRepo.MarkOperatorDeleted(ctx, userID)
	if err != nil {
		log.Error("Failed to delete user from operators: ", err)
		return err
	}

	log.Info("Remove operator for tasks",
		slog.Any("task_ids", task_ids),
		slog.Int64("user_id", userID))

	return nil
}
