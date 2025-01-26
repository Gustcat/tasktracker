package user

import (
	"github.com/Gustcat/auth/internal/client/db"
	"github.com/Gustcat/auth/internal/repository"
	"github.com/Gustcat/auth/internal/service"
)

type serv struct {
	userRepository repository.UserRepository
	txManager      db.TxManager
}

func NewServ(userRepository repository.UserRepository, txManager db.TxManager) service.UserService {
	return &serv{
		userRepository: userRepository,
		txManager:      txManager,
	}
}
