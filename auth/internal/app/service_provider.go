package app

import (
	"context"
	accessAPI "github.com/Gustcat/auth/internal/api/access"
	authAPI "github.com/Gustcat/auth/internal/api/auth"
	userAPI "github.com/Gustcat/auth/internal/api/user"
	"github.com/Gustcat/auth/internal/client/db"
	"github.com/Gustcat/auth/internal/client/db/pg"
	"github.com/Gustcat/auth/internal/client/db/transaction"
	"github.com/Gustcat/auth/internal/closer"
	"github.com/Gustcat/auth/internal/config"
	"github.com/Gustcat/auth/internal/config/env"
	"github.com/Gustcat/auth/internal/repository"
	accessRepository "github.com/Gustcat/auth/internal/repository/access"
	userRepository "github.com/Gustcat/auth/internal/repository/user"
	"github.com/Gustcat/auth/internal/service"
	accessService "github.com/Gustcat/auth/internal/service/access"
	authService "github.com/Gustcat/auth/internal/service/auth"
	userService "github.com/Gustcat/auth/internal/service/user"
	"log"
)

type serviceProvider struct {
	pgConfig      config.PGConfig
	grpcConfig    config.GRPCConfig
	httpConfig    config.HTTPConfig
	swaggerConfig config.SwaggerConfig
	tokenConfig   config.TokenConfig

	dbClient  db.Client
	txManager db.TxManager

	userRepository   repository.UserRepository
	accessRepository repository.AccessRepository

	userService   service.UserService
	authService   service.AuthService
	accessService service.AccessService

	userImpl   *userAPI.Implementation
	authImpl   *authAPI.Implementation
	accessImpl *accessAPI.Implementation
}

func newServiceProvider() *serviceProvider { return &serviceProvider{} }

func (sp *serviceProvider) PGConfig() config.PGConfig {
	if sp.pgConfig == nil {
		cfg, err := env.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}
		sp.pgConfig = cfg
	}

	return sp.pgConfig
}

func (sp *serviceProvider) GRPCConfig() config.GRPCConfig {
	if sp.grpcConfig == nil {
		cfg, err := env.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}
		sp.grpcConfig = cfg
	}

	return sp.grpcConfig
}

func (sp *serviceProvider) HTTPConfig() config.HTTPConfig {
	if sp.httpConfig == nil {
		cfg, err := env.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get http config: %s", err.Error())
		}

		sp.httpConfig = cfg
	}

	return sp.httpConfig
}

func (sp *serviceProvider) SwaggerConfig() config.SwaggerConfig {
	if sp.swaggerConfig == nil {
		cfg, err := env.NewSwaggerConfig()
		if err != nil {
			log.Fatalf("failed to get swagger config: %s", err.Error())
		}

		sp.swaggerConfig = cfg
	}

	return sp.swaggerConfig
}

func (sp *serviceProvider) TokenConfig() config.TokenConfig {
	if sp.tokenConfig == nil {
		cfg, err := env.NewTokenConfig()
		if err != nil {
			log.Fatalf("failed to get swagger config: %s", err.Error())
		}

		sp.tokenConfig = cfg
	}

	return sp.tokenConfig
}

func (sp *serviceProvider) DBClient(ctx context.Context) db.Client {
	if sp.dbClient == nil {
		client, err := pg.New(ctx, sp.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to connect to db: %s", err.Error())
		}

		err = client.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("failed to ping pgxpool: %s", err.Error())
		}
		closer.Add(client.Close)

		sp.dbClient = client
	}

	return sp.dbClient
}

func (sp *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if sp.txManager == nil {
		sp.txManager = transaction.NewTransactionManager(sp.DBClient(ctx).DB())
	}

	return sp.txManager
}

func (sp *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if sp.userRepository == nil {
		sp.userRepository = userRepository.NewRepository(sp.DBClient(ctx))
	}

	return sp.userRepository
}

func (sp *serviceProvider) AccessRepository(ctx context.Context) repository.AccessRepository {
	if sp.accessRepository == nil {
		sp.accessRepository = accessRepository.NewRepository(sp.DBClient(ctx))
	}

	return sp.accessRepository
}

func (sp *serviceProvider) UserService(ctx context.Context) service.UserService {
	if sp.userService == nil {
		sp.userService = userService.NewServ(sp.UserRepository(ctx), sp.TxManager(ctx))
	}

	return sp.userService
}

func (sp *serviceProvider) AccessService(ctx context.Context) service.AccessService {
	if sp.accessService == nil {
		sp.accessService = accessService.NewService(sp.AccessRepository(ctx), sp.TokenConfig())
	}

	return sp.accessService
}

func (sp *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if sp.authService == nil {
		sp.authService = authService.NewService(sp.UserRepository(ctx), sp.TokenConfig())
	}

	return sp.authService
}

func (sp *serviceProvider) UserImpl(ctx context.Context) *userAPI.Implementation {
	if sp.userImpl == nil {
		sp.userImpl = userAPI.NewImplementation(sp.UserService(ctx))
	}

	return sp.userImpl
}

func (sp *serviceProvider) AuthImpl(ctx context.Context) *authAPI.Implementation {
	if sp.authImpl == nil {
		sp.authImpl = authAPI.NewImplementation(sp.AuthService(ctx))
	}

	return sp.authImpl
}

func (sp *serviceProvider) AccessImpl(ctx context.Context) *accessAPI.Implementation {
	if sp.accessImpl == nil {
		sp.accessImpl = accessAPI.NewImplementation(sp.AccessService(ctx), sp.TokenConfig())
	}

	return sp.accessImpl
}
