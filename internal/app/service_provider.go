package app

import (
	"context"
	userAPI "github.com/Gustcat/auth/internal/api"
	"github.com/Gustcat/auth/internal/client/db"
	"github.com/Gustcat/auth/internal/client/db/pg"
	"github.com/Gustcat/auth/internal/client/db/transaction"
	"github.com/Gustcat/auth/internal/closer"
	"github.com/Gustcat/auth/internal/config"
	"github.com/Gustcat/auth/internal/config/env"
	"github.com/Gustcat/auth/internal/repository"
	userRepository "github.com/Gustcat/auth/internal/repository/user"
	"github.com/Gustcat/auth/internal/service"
	userService "github.com/Gustcat/auth/internal/service/user"
	"log"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	dbClient       db.Client
	txManager      db.TxManager
	userRepository repository.UserRepository
	userService    service.UserService
	userImpl       *userAPI.Implementation
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

func (sp *serviceProvider) UserService(ctx context.Context) service.UserService {
	if sp.userService == nil {
		sp.userService = userService.NewServ(sp.UserRepository(ctx), sp.TxManager(ctx))
	}

	return sp.userService
}

func (sp *serviceProvider) UserImpl(ctx context.Context) *userAPI.Implementation {
	if sp.userImpl == nil {
		sp.userImpl = userAPI.NewImplementation(sp.UserService(ctx))
	}

	return sp.userImpl
}
