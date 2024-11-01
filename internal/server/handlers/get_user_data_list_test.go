package handlers

import (
	"context"
	"github.com/SversusN/keeper/internal/utils"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/SversusN/keeper/internal/server/storage"
	mock_storage "github.com/SversusN/keeper/internal/server/storage/mocks"
	pb "github.com/SversusN/keeper/pkg/grpc"
)

func TestServer_GetUserDataList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB = mock_storage.NewMockrepository(ctrl)
	srv.Storage = mockDB
	mockDB.EXPECT().GetUserData(gomock.Any(), int64(3)).Return(nil, storage.ErrGetUserData).AnyTimes()
	mockDB.EXPECT().GetUserData(gomock.Any(), gomock.Any()).Return([]storage.ShortRecord{{
		Name:     "testName",
		DataType: "password",
		ID:       1,
		Version:  1,
	}}, nil).AnyTimes()

	type args struct {
		in *pb.UserDataListRequest
	}
	tests := []struct {
		name    string
		args    args
		userID  int64
		want    *pb.UserDataListResponse
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				in: &pb.UserDataListRequest{},
			},
			userID: 1,
			want: &pb.UserDataListResponse{
				Data: []*pb.UserDataNested{{
					Id:       1,
					Name:     "testName",
					DataType: "password",
					Version:  1,
				}},
			},
			wantErr: false,
		},
		{
			name: "internal error",
			args: args{
				in: &pb.UserDataListRequest{},
			},
			userID:  3,
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := srv
			ctx = context.WithValue(context.Background(), utils.UserIDContextKey, tt.userID)
			got, err := s.GetUserDataList(ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserDataList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserDataList() got = %v, want %v", got, tt.want)
			}
		})
	}
}
