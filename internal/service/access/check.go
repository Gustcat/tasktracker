package access

import (
	"context"
	"errors"
	"github.com/Gustcat/auth/internal/utils"
)

func (s *serv) Check(ctx context.Context, endpointAddress string, accessToken string) error {
	claims, err := utils.VerifyToken(accessToken, []byte(s.tokenConfig.AccessTokenSecretKey()))
	if err != nil {
		return errors.New("invalid access token")
	}

	err = s.accessRepository.Check(ctx, claims.Role, endpointAddress)
	if err != nil {
		return errors.New("endpoint is not available to the user")
	}

	return nil
}
