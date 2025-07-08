package ctxutils

import (
	"context"
	"github.com/Gustcat/task-server/internal/model"
)

type ctxKey string

const (
	LoggerKey ctxKey = "logger"
	UserKey   ctxKey = "user"
)

func UserFromContext(ctx context.Context) *model.User {
	user, ok := ctx.Value(UserKey).(*model.User)
	if !ok {
		return nil
	}
	return user
}
