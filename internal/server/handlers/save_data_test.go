package handlers

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	mock_storage "github.com/SversusN/keeper/internal/server/storage/mocks"
	pb "github.com/SversusN/keeper/pkg/grpc"
)

func TestServer_SaveData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB = mock_storage.NewMockrepository(ctrl)
	srv.Storage = mockDB
	mockDB.EXPECT().SaveUserData(gomock.Any(), gomock.Any(), "error", gomock.Any(), gomock.Any()).Return(
		errors.ErrSaveUserData).AnyTimes()
	mockDB.EXPECT().SaveUserData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
		nil).AnyTimes()

	type args struct {
		in *pb.SaveDataRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.SaveDataResponse
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				in: &pb.SaveDataRequest{
					Name:     "test",
					Data:     []byte("test"),
					DataType: "password",
				},
			},
			want:    &pb.SaveDataResponse{Result: "OK"},
			wantErr: false,
		},
		{
			name: "internal error",
			args: args{
				in: &pb.SaveDataRequest{
					Name:     "error",
					Data:     []byte("test"),
					DataType: "password",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := srv
			got, err := s.SaveData(ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SaveData() got = %v, want %v", got, tt.want)
			}
		})
	}
}
