package app

import (
	"context"
	userAPI "github.com/Gustcat/auth/internal/api"
	"github.com/Gustcat/auth/internal/closer"
	"github.com/Gustcat/auth/internal/config"
	"github.com/Gustcat/auth/internal/config/env"
	"github.com/Gustcat/auth/internal/repository"
	userRepository "github.com/Gustcat/auth/internal/repository/user"
	"github.com/Gustcat/auth/internal/service"
	userService "github.com/Gustcat/auth/internal/service/user"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	pgPool         *pgxpool.Pool
	userRepository repository.UserRepository
	userService    service.UserService
	userImpl       *userAPI.Implementation
}

func NewServiceProvider() *serviceProvider { return &serviceProvider{} }

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

func (sp *serviceProvider) PgPool(ctx context.Context) *pgxpool.Pool {
	if sp.pgPool == nil {
		pool, err := pgxpool.Connect(ctx, sp.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to connect to pgxpool: %s", err.Error())
		}

		err = pool.Ping(ctx)
		if err != nil {
			log.Fatalf("failed to ping pgxpool: %s", err.Error())
		}
		closer.Add(func() error {
			pool.Close()
			return nil
		})
		sp.pgPool = pool
	}

	return sp.pgPool
}

func (sp *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if sp.userRepository == nil {
		sp.userRepository = userRepository.NewRepository(sp.PgPool(ctx))
	}

	return sp.userRepository
}

func (sp *serviceProvider) UserService(ctx context.Context) service.UserService {
	if sp.userService == nil {
		sp.userService = userService.NewServ(sp.UserRepository(ctx))
	}

	return sp.userService
}

func (sp *serviceProvider) UserImpl(ctx context.Context) *userAPI.Implementation {
	if sp.userImpl == nil {
		sp.userImpl = userAPI.NewImplementation(sp.UserService(ctx))
	}

	return sp.userImpl
}
