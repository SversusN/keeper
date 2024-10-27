package handlers

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/SversusN/keeper/internal/server/storage"
	mock_storage "github.com/SversusN/keeper/internal/server/storage/mocks"
	pb "github.com/SversusN/keeper/pkg/grpc"
)

func TestServer_GetUserData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB = mock_storage.NewMockrepository(ctrl)
	srv.Storage = mockDB
	mockDB.EXPECT().FindUserRecord(gomock.Any(), int64(0), gomock.Any()).Return(nil, storage.ErrFindUserRecord).AnyTimes()
	mockDB.EXPECT().FindUserRecord(gomock.Any(), int64(2), gomock.Any()).Return(nil, storage.ErrNowRows).AnyTimes()
	mockDB.EXPECT().FindUserRecord(gomock.Any(), gomock.Any(), gomock.Any()).Return(testRecord, nil).AnyTimes()

	type args struct {
		in *pb.UserDataRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.UserDataResponse
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				in: &pb.UserDataRequest{Id: 1},
			},
			want: &pb.UserDataResponse{
				Id:       1,
				Name:     "testName",
				Data:     []byte("test"),
				DataType: "password",
				Version:  1,
				CreateAt: "",
			},
			wantErr: false,
		},
		{
			name: "internal error",
			args: args{
				in: &pb.UserDataRequest{Id: 0},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "no content",
			args: args{
				in: &pb.UserDataRequest{Id: 2},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := srv
			got, err := s.GetUserData(ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserData() got = %v, want %v", got, tt.want)
			}
		})
	}
}
