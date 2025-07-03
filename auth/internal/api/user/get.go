package user

import (
	"context"
	"errors"
	"github.com/Gustcat/auth/internal/converter"
	"github.com/Gustcat/auth/internal/logger"
	"github.com/Gustcat/auth/internal/repository"
	desc "github.com/Gustcat/auth/pkg/user_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	logger.Info("Getting user", zap.Int64("user_id", req.GetId()))

	id, userinfo, createdAt, updatedAt, err := i.userService.Get(ctx, req.GetId())
	if errors.Is(err, repository.ErrUserNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if err != nil {
		return nil, err
	}

	logger.Info("Got user", zap.Int64("id", id))

	return &desc.GetResponse{
		Id:        id,
		Info:      converter.ToUserInfoFromService(*userinfo),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}
