package access

import (
	"github.com/Gustcat/auth/internal/config"
	"github.com/Gustcat/auth/internal/repository"
	"github.com/Gustcat/auth/internal/service"
)

type serv struct {
	accessRepository repository.AccessRepository
	tokenConfig      config.TokenConfig
}

func NewService(accessRepository repository.AccessRepository, tokenConfig config.TokenConfig) service.AccessService {
	return &serv{
		accessRepository: accessRepository,
		tokenConfig:      tokenConfig,
	}
}
