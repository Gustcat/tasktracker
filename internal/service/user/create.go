package user

import (
	"context"
	"github.com/Gustcat/auth/internal/model"
)

func (s *serv) Create(ctx context.Context, info *model.UserInfo, pwd string) (int64, error) {
	id, err := s.userRepository.Create(ctx, info, pwd)
	if err != nil {
		return 0, err
	}

	return id, nil
}
