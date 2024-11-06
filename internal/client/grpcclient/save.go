package grpcclient

import (
	"context"
	"fmt"
	"github.com/SversusN/keeper/internal/utils/encrypter"

	"github.com/SversusN/keeper/internal/client/models"
	pb "github.com/SversusN/keeper/pkg/grpc"
)

// SaveUserData – метод сохранения данных пользователя на сервер.
func (c *Client) SaveUserData(model *models.UserData) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	req := &pb.SaveDataRequest{
		Name:     model.Name,
		Data:     encrypter.Encrypt(model.Data, c.config.PassPhrase),
		DataType: model.DataType,
	}
	_, err := c.gRPCClient.SaveData(ctx, req)
	if err != nil {
		return fmt.Errorf("%w: gRPC SaveUserData error: %w", ErrRequest, err)
	}

	return nil
}
