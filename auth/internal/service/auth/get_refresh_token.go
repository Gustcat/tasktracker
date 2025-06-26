package auth

import (
	"context"
	"github.com/Gustcat/auth/internal/model"
	"github.com/Gustcat/auth/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serv) GetRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	claims, err := utils.VerifyToken(refreshToken, []byte(s.tokenConfig.RefreshTokenSecretKey()))
	if err != nil {
		return "", status.Errorf(codes.Aborted, "invalid refresh token")
	}

	newRefreshToken, err := utils.GenerateToken(model.UserInfo{
		Name: claims.Username,
		Role: claims.Role,
	},
		[]byte(s.tokenConfig.RefreshTokenSecretKey()),
		s.tokenConfig.RefreshTokenExpiration(),
	)
	if err != nil {
		return "", err
	}

	return newRefreshToken, nil
}
