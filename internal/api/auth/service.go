package auth

import (
	"github.com/Gustcat/auth/internal/service"
	descAuth "github.com/Gustcat/auth/pkg/auth_v1"
)

type Implementation struct {
	descAuth.UnimplementedAuthV1Server
	authService service.AuthService
}

func NewImplementation(authService service.AuthService) *Implementation {
	return &Implementation{authService: authService}
}
