package grpcclient

import (
	"errors"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/SversusN/keeper/internal/client/config"
	"github.com/SversusN/keeper/internal/client/interceptors"
	pb "github.com/SversusN/keeper/pkg/grpc"
)

// Client – структура gRPC клиента для общения с сервером.
type Client struct {
	gRPCClient pb.KeeperClient
	config     *config.Config
	authToken  string
	timeout    time.Duration
}

var ErrRequest = errors.New(`request error`)

// NewGRPCClient – функция инициализации gRPC клиента.
func NewGRPCClient(c *config.Config) (*Client, error) {
	client := &Client{
		config:  c,
		timeout: time.Duration(c.ConnectionTimeout) * time.Second,
	}
	conn, err := grpc.NewClient(
		c.Host,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(interceptors.AuthInterceptor(client)),
	)
	if err != nil {
		return nil, fmt.Errorf("gRPC connection refused: %w", err)
	}
	gRPCClient := pb.NewKeeperClient(conn)
	client.gRPCClient = gRPCClient

	return client, nil
}
