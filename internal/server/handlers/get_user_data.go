package handlers

import (
	"context"
	"errors"
	"github.com/SversusN/keeper/internal/utils"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/SversusN/keeper/internal/server/storage"
	pb "github.com/SversusN/keeper/pkg/grpc"
)

// GetUserData – метод получения сохранённых данных пользователя.
func (s *Server) GetUserData(ctx context.Context, in *pb.UserDataRequest) (*pb.UserDataResponse, error) {
	userID, ok := ctx.Value(utils.UserIDContextKey).(int64)
	if !ok {
		s.Logger.Log.Error(missingKeyErrText)
		return nil, status.Error(codes.Internal, http.StatusText(http.StatusInternalServerError))
	}

	record, err := s.Storage.FindUserRecord(ctx, in.Id, userID)
	if err != nil {
		if errors.Is(err, storage.ErrNowRows) {
			return nil, status.Error(codes.NotFound, http.StatusText(http.StatusNoContent))
		}
		s.Logger.Log.Error(err)
		return nil, status.Error(codes.Internal, http.StatusText(http.StatusInternalServerError))
	}

	return &pb.UserDataResponse{
		Id:       record.ID,
		Name:     record.Name,
		Data:     record.Data,
		DataType: record.DataType,
		Version:  record.Version,
		CreateAt: record.CreatedAt,
	}, nil
}
