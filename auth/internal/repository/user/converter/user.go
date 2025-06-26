package converter

import (
	"database/sql"
	"github.com/Gustcat/auth/internal/model"
	modelRepo "github.com/Gustcat/auth/internal/repository/user/model"
	"time"
)

func ToUserFromRepo(user *modelRepo.User) (int64, *model.UserInfo, time.Time, sql.NullTime) {
	return user.ID, ToUserInfoFromRepo(user.Info), user.CreatedAt, user.UpdatedAt
}

func ToUserInfoFromRepo(userinfo modelRepo.UserInfo) *model.UserInfo {
	return &model.UserInfo{
		Name:  userinfo.Name,
		Email: userinfo.Email,
		Role:  userinfo.Role,
	}
}
