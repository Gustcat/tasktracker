package repository

import (
	"context"
	"database/sql"
	"github.com/Gustcat/auth/internal/model"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, info *model.UserInfo, pwd string) (int64, error)
	Get(ctx context.Context, id int64) (int64, *model.UserInfo, time.Time, sql.NullTime, error)
	Update(ctx context.Context, id int64, name string, email string) error
	Delete(ctx context.Context, id int64) error
}
