package auth

import (
	"context"
	"github.com/Gustcat/auth/internal/model"
	"github.com/Gustcat/auth/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serv) GetAccessToken(ctx context.Context, refreshToken string) (string, error) {
	claims, err := utils.VerifyToken(refreshToken, []byte(s.tokenConfig.RefreshTokenSecretKey()))
	if err != nil {
		return "", status.Errorf(codes.Aborted, "invalid refresh token")
	}

	newAccessToken, err := utils.GenerateToken(model.UserToken{
		Name: claims.Name,
		Role: claims.Role,
		ID:   claims.ID,
	},
		[]byte(s.tokenConfig.AccessTokenSecretKey()),
		s.tokenConfig.AccessTokenExpiration(),
	)
	if err != nil {
		return "", err
	}

	return newAccessToken, nil
}
