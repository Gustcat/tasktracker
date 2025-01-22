package api

import (
	"context"
	desc "github.com/Gustcat/auth/pkg/user_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	err := i.userService.Delete(ctx, req.GetId())
	if err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "не удалось выполнить SQL-запрос: %v", err)
	}

	return &emptypb.Empty{}, nil
}
