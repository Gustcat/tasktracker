package service

import (
	"context"
	"github.com/Gustcat/auth/internal/model"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserService interface {
	Create(ctx context.Context, info *model.UserInfo, pwd string) (int64, error)
	Get(ctx context.Context, id int64) (int64, *model.UserInfo, *timestamppb.Timestamp, *timestamppb.Timestamp, error)
	Update(ctx context.Context, id int64, name string, email string) error
	Delete(ctx context.Context, id int64) error
}
