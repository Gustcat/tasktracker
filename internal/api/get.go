package api

import (
	"context"
	"github.com/Gustcat/auth/internal/converter"
	desc "github.com/Gustcat/auth/pkg/user_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func (i *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("fetching user with id: %d", req.Id)

	id, userinfo, createdAt, updatedAt, err := i.userService.Get(ctx, req.GetId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &desc.GetResponse{
		Id:        id,
		Info:      converter.ToUserInfoFromService(*userinfo),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}
