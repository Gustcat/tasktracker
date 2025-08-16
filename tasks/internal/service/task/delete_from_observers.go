package task

import (
	"context"
	"github.com/Gustcat/task-server/internal/logger"
	"log/slog"
)

func (s *Serv) DeleteUserFromObservers(ctx context.Context, userID int64) error {
	const op = "service.task.DeleteUserFromObservers"

	log := logger.LogFromContextAddOP(ctx, op)
	log.Info("DeleteUserFromObservers called")

	// TODO: select-запрос на получение task_id для логирования, обернуть в транзакцию

	err := s.watcherRepo.DeleteUser(ctx, userID)
	if err != nil {
		log.Error("Failed to delete user from observers: ", err)
		return err
	}

	log.Info("Remove watcher for tasks",
		//slog.Any("task_ids", taskIDs),
		slog.Int64("user_id", userID))
	return nil
}
