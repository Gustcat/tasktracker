package auth

import (
	"github.com/Gustcat/auth/internal/config"
	"github.com/Gustcat/auth/internal/repository"
	"github.com/Gustcat/auth/internal/service"
)

type serv struct {
	userRepository repository.UserRepository
	tokenConfig    config.TokenConfig
}

func NewService(userRepository repository.UserRepository, tokenConfig config.TokenConfig) service.AuthService {
	return &serv{
		userRepository: userRepository,
		tokenConfig:    tokenConfig,
	}
}
