package grpcclient

import (
	"context"
	"fmt"
	"github.com/SversusN/keeper/internal/utils/encrypter"

	"github.com/SversusN/keeper/internal/client/models"
	pb "github.com/SversusN/keeper/pkg/grpc"
)

// UpdateUserData – обновление данных пользователя на сервере.
func (c *Client) UpdateUserData(model *models.UserData) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	req := &pb.UpdateUserDataRequest{
		Id:      model.ID,
		Data:    encrypter.Encrypt(model.Data, c.config.PassPhrase),
		Version: model.Version,
	}
	_, err := c.gRPCClient.UpdateUserData(ctx, req)
	if err != nil {
		return fmt.Errorf("%w: gRPC UpdataUserData error: %w", ErrRequest, err)
	}

	return nil
}
