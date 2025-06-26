package access

import (
	"github.com/Gustcat/auth/internal/config"
	"github.com/Gustcat/auth/internal/service"
	descAccess "github.com/Gustcat/auth/pkg/access_v1"
)

type Implementation struct {
	descAccess.UnimplementedAccessV1Server
	accessService service.AccessService
	tokenConfig   config.TokenConfig
}

func NewImplementation(accessService service.AccessService, tokenConfig config.TokenConfig) *Implementation {
	return &Implementation{
		accessService: accessService,
		tokenConfig:   tokenConfig,
	}
}
