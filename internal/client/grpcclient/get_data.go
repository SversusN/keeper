package grpcclient

import (
	"context"
	"fmt"
	"github.com/SversusN/keeper/internal/utils/encrypter"

	"github.com/SversusN/keeper/internal/client/models"
	pb "github.com/SversusN/keeper/pkg/grpc"
)

// GetUserData – метод получения данных пользователя с сервера.
func (c *Client) GetUserData(model models.UserDataModel) (*models.UserData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	req := &pb.UserDataRequest{Id: model.ID}
	res, err := c.gRPCClient.GetUserData(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("%w: gRPC GetUserData error: %w", ErrRequest, err)
	}

	return &models.UserData{
		ID:       res.Id,
		Name:     res.Name,
		DataType: res.DataType,
		Version:  res.Version,
		Data:     encrypter.Decrypt(res.Data, c.config.PassPhrase),
	}, nil
}

// GetUserDataList – метод получения всех сохранённых данных (мета-данных) пользователя с сервера.
func (c *Client) GetUserDataList() ([]models.UserDataList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	req := &pb.UserDataListRequest{}
	res, err := c.gRPCClient.GetUserDataList(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("%w: gRPC GetUserDataList error: %w", ErrRequest, err)
	}

	records := make([]models.UserDataList, 0, len(res.Data))
	for _, el := range res.Data {
		rec := models.UserDataList{
			ID:       el.Id,
			Name:     el.Name,
			DataType: el.DataType,
			Version:  el.Version,
		}
		records = append(records, rec)
	}

	return records, nil
}
