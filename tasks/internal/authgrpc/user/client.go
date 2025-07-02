package user

import (
	"context"
	"fmt"
	descUser "github.com/Gustcat/auth/pkg/user_v1"
	"github.com/Gustcat/task-server/internal/model"
	"github.com/Gustcat/task-server/internal/service"
)

type Client struct {
	grpc descUser.UserV1Client
}

func NewClient(grpc descUser.UserV1Client) service.AuthService {
	return &Client{grpc: grpc}
}

func (c *Client) GetUser(ctx context.Context, id int64) (*model.User, error) {
	resp, err := c.grpc.Get(ctx, &descUser.GetRequest{Id: id})
	if err != nil {
		return nil, fmt.Errorf("GetUserByID: %w", err)
	}

	return &model.User{
		ID:    resp.Id,
		Name:  resp.Info.Name,
		Email: resp.Info.Email,
		Role:  model.Role(resp.Info.Role),
	}, nil
}
