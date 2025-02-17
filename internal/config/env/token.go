package env

import (
	"errors"
	"github.com/Gustcat/auth/internal/config"
	"os"
	"strconv"
	"time"
)

var _ config.TokenConfig = (*tokenConfig)(nil)

const (
	accessTokenSecretKeyEnvName   = "ACCESS_TOKEN_SECRET"
	accessTokenExpirationEnvName  = "ACCESS_TOKEN_EXPIRATION"
	refreshTokenSecretKeyEnvName  = "REFRESH_TOKEN_SECRET"
	refreshTokenExpirationEnvName = "REFRESH_TOKEN_EXPIRATION"
	authPrefixEnvName             = "AUTH_PREFIX"
)

type tokenConfig struct {
	accessTokenSecretKey         string
	accessTokenExpirationMinute  int64
	refreshTokenSecretKey        string
	refreshTokenExpirationMinute int64
	authPrefix                   string
}

func NewTokenConfig() (*tokenConfig, error) {
	accessTokenSecretKey := os.Getenv(accessTokenSecretKeyEnvName)
	if len(accessTokenSecretKey) == 0 {
		return nil, errors.New("accessTokenSecretKey not found")
	}

	strAccessTokenExpiration := os.Getenv(accessTokenExpirationEnvName)
	if len(strAccessTokenExpiration) == 0 {
		return nil, errors.New("accessTokenExpiration not found")
	}
	accessTokenExpirationMinute, err := strconv.ParseInt(strAccessTokenExpiration, 10, 64)
	if err != nil {
		return nil, errors.New("accessTokenExpiration is not int")
	}

	refreshTokenSecretKey := os.Getenv(refreshTokenSecretKeyEnvName)
	if len(refreshTokenSecretKey) == 0 {
		return nil, errors.New("refreshTokenSecretKey not found")
	}

	strRefreshTokenExpiration := os.Getenv(refreshTokenExpirationEnvName)
	if len(strRefreshTokenExpiration) == 0 {
		return nil, errors.New("refreshTokenExpiration not found")
	}
	refreshTokenExpirationMinute, err := strconv.ParseInt(strRefreshTokenExpiration, 10, 64)
	if err != nil {
		return nil, errors.New("accessTokenExpiration is not int")
	}

	authPrefix := os.Getenv(authPrefixEnvName)
	if len(authPrefix) == 0 {
		return nil, errors.New("authPrefix not found")
	}
	return &tokenConfig{
		accessTokenSecretKey:         accessTokenSecretKey,
		accessTokenExpirationMinute:  accessTokenExpirationMinute,
		refreshTokenExpirationMinute: refreshTokenExpirationMinute,
		refreshTokenSecretKey:        refreshTokenSecretKey,
		authPrefix:                   authPrefix,
	}, nil
}

func (cfg *tokenConfig) AccessTokenSecretKey() string {
	return cfg.accessTokenSecretKey
}

func (cfg *tokenConfig) AccessTokenExpiration() time.Duration {
	return time.Duration(cfg.accessTokenExpirationMinute) * time.Minute
}

func (cfg *tokenConfig) RefreshTokenSecretKey() string {
	return cfg.refreshTokenSecretKey
}

func (cfg *tokenConfig) RefreshTokenExpiration() time.Duration {
	return time.Duration(cfg.refreshTokenExpirationMinute) * time.Minute
}

func (cfg *tokenConfig) AuthPrefix() string {
	return cfg.authPrefix
}
