package task

import (
	"context"
	"errors"
	"github.com/Gustcat/task-server/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

func (s *Serv) validateUser(ctx context.Context, userId int64) (*model.User, error) {
	user, err := s.authClient.GetUser(ctx, userId)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}
