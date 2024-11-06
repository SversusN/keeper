package grpcclient

import (
	"context"
	"fmt"

	"github.com/SversusN/keeper/internal/client/models"
	pb "github.com/SversusN/keeper/pkg/grpc"
)

// GetAuthToken – метод получения AuthToken пользователя.
func (c *Client) GetAuthToken() string {
	return c.authToken
}

// Register – метод регистрации пользователя на сервере.
func (c *Client) Register(model models.AuthModel) (models.AuthToken, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	req := &pb.RegisterRequest{
		Login:    model.Login,
		Password: model.Password,
	}
	res, err := c.gRPCClient.Register(ctx, req)
	if err != nil {
		return "", fmt.Errorf("%w: gRPC Register error: %w", ErrRequest, err)
	}
	c.authToken = res.Token

	return models.AuthToken(res.Token), nil
}

// SignIn – метод авторизации клиента на сервере.
func (c *Client) SignIn(model models.AuthModel) (models.AuthToken, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	req := &pb.SignInRequest{
		Login:    model.Login,
		Password: model.Password,
	}
	res, err := c.gRPCClient.SignIn(ctx, req)
	if err != nil {
		return "", fmt.Errorf("%w: gRPC SignIn error: %w", ErrRequest, err)
	}
	c.authToken = res.Token
	return models.AuthToken(res.Token), nil
}
