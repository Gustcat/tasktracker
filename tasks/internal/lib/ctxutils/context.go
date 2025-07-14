package ctxutils

import (
	"context"
	"errors"
	"github.com/Gustcat/task-server/internal/model"
)

type ctxKey string

const (
	LoggerKey ctxKey = "logger"
	UserKey   ctxKey = "user"
	TxKey     ctxKey = "tx"
)

var (
	ErrCurrentUserNotFound = errors.New("current user not found in context")
)

func UserFromContext(ctx context.Context) (*model.User, error) {
	user, ok := ctx.Value(UserKey).(*model.User)
	if !ok {
		return nil, ErrCurrentUserNotFound
	}

	return user, nil
}
