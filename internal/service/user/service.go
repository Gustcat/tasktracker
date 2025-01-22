package user

import (
	"github.com/Gustcat/auth/internal/repository"
	"github.com/Gustcat/auth/internal/service"
)

type serv struct {
	userRepository repository.UserRepository
}

func NewServ(userRepository repository.UserRepository) service.UserService {
	return &serv{
		userRepository: userRepository,
	}
}
