package user

import (
	"context"
	"github.com/Gustcat/auth/internal/logger"
	"github.com/Gustcat/auth/internal/model"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *serv) Get(ctx context.Context, id int64) (int64, *model.UserInfo, *timestamppb.Timestamp, *timestamppb.Timestamp, error) {
	logger.Info("Getting user...", zap.Int64("id", id))

	id, userinfo, createdAt, updatedAtTime, err := s.userRepository.Get(ctx, id)
	if err != nil {
		return 0, nil, &timestamppb.Timestamp{}, &timestamppb.Timestamp{}, err
	}
	var updatedAt *timestamppb.Timestamp
	if updatedAtTime.Valid {
		updatedAt = timestamppb.New(updatedAtTime.Time)
	}

	logger.Info("Got user", zap.Int64("id", id))

	return id, userinfo, timestamppb.New(createdAt), updatedAt, nil
}
