package user

import (
	"context"
	"github.com/Gustcat/auth/internal/logger"
	"go.uber.org/zap"
)

func (s *serv) Delete(ctx context.Context, id int64) error {
	logger.Info("Deleting user...", zap.Int64("id", id))
	err := s.userRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	logger.Info("User deleted", zap.Int64("id", id))

	return nil
}
