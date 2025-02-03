package api

import (
	"context"
	desc "github.com/Gustcat/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {

	err := i.userService.Update(ctx, req.GetId(), req.GetName(), req.GetEmail())
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
