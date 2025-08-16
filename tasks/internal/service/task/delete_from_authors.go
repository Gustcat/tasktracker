package task

import (
	"context"
	"github.com/Gustcat/task-server/internal/logger"
	"log/slog"
)

func (s *Serv) DeleteUserFromAuthors(ctx context.Context, userID int64) error {
	const op = "service.task.DeleteUserFromAuthors"

	log := logger.LogFromContextAddOP(ctx, op)
	log.Info("DeleteUserFromAuthors called")

	task_ids, err := s.taskRepo.MarkAuthorDeleted(ctx, userID)
	if err != nil {
		log.Error("Failed to delete user from authors: ", err)
		return err
	}

	log.Info("Remove author for tasks",
		slog.Any("task_ids", task_ids),
		slog.Int64("user_id", userID))

	return nil
}
