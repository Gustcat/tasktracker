package user

import (
	"github.com/Gustcat/auth/internal/client/db"
	"github.com/Gustcat/auth/internal/repository"
	"github.com/Gustcat/auth/internal/service"
	"github.com/Gustcat/shared-lib/kafka_common"
)

type serv struct {
	userRepository repository.UserRepository
	txManager      db.TxManager
	producer       *kafka_common.Producer
}

func NewServ(userRepository repository.UserRepository, txManager db.TxManager, producer *kafka_common.Producer) service.UserService {
	return &serv{
		userRepository: userRepository,
		txManager:      txManager,
		producer:       producer,
	}
}

func NewMockService(deps ...interface{}) service.UserService {
	srv := serv{}

	for _, v := range deps {
		switch s := v.(type) {
		case repository.UserRepository:
			srv.userRepository = s
		}
	}

	return &srv
}
