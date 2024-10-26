package handlers

import (
	"context"
	"errors"
	errors2 "github.com/SversusN/keeper/internal/server/internalerrors"
	"github.com/SversusN/keeper/internal/utils"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
		if errors.Is(err, errors2.ErrNoRows) {
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

// GetUserDataList Коллекция данных пользователя
func (s *Server) GetUserDataList(ctx context.Context, in *pb.UserDataListRequest) (*pb.UserDataListResponse, error) {
	userID, ok := ctx.Value(utils.UserIDContextKey).(int64)
	if !ok {
		s.Logger.Log.Error(missingKeyErrText)
		return nil, status.Error(codes.Internal, http.StatusText(http.StatusInternalServerError))
	}

	userRecords, err := s.Storage.GetUserData(ctx, userID)
	if err != nil {
		s.Logger.Log.Error(err)
		return nil, status.Error(codes.Internal, http.StatusText(http.StatusInternalServerError))
	}

	records := make([]*pb.UserDataNested, 0, len(userRecords))
	for _, rec := range userRecords {
		data := &pb.UserDataNested{
			Id:       rec.ID,
			Name:     rec.Name,
			DataType: rec.DataType,
			Version:  rec.Version,
			CreateAt: rec.CreatedAt,
		}
		records = append(records, data)
	}

	return &pb.UserDataListResponse{
		Data: records,
	}, nil
}
