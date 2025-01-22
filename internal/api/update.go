package api

import (
	"context"
	desc "github.com/Gustcat/auth/pkg/user_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {

	err := i.userService.Update(ctx, req.GetId(), req.GetName(), req.GetEmail())
	if err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "не удалось выполнить SQL-запрос: %v", err)
	}
	return &emptypb.Empty{}, nil
}
