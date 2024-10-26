package handlers

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	mock_storage "github.com/SversusN/keeper/internal/server/storage/mocks"
	pb "github.com/SversusN/keeper/pkg/grpc"
)

func TestServer_UpdateUserData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB = mock_storage.NewMockrepository(ctrl)
	srv.Storage = mockDB
	mockDB.EXPECT().FindUserRecord(gomock.Any(), int64(0), gomock.Any()).Return(nil, errors.ErrFindUserRecord).AnyTimes()
	mockDB.EXPECT().FindUserRecord(gomock.Any(), int64(2), gomock.Any()).Return(nil, errors.ErrNoRows).AnyTimes()
	mockDB.EXPECT().FindUserRecord(gomock.Any(), gomock.Any(), gomock.Any()).Return(testRecord, nil).AnyTimes()
	mockDB.EXPECT().UpdateUserRecord(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	type args struct {
		in *pb.UpdateUserDataRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.UpdateUserDataResponse
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				in: &pb.UpdateUserDataRequest{
					Id:      1,
					Data:    []byte("test"),
					Version: 1,
				},
			},
			want: &pb.UpdateUserDataResponse{
				Result: "OK",
			},
			wantErr: false,
		},
		{
			name: "invalid version",
			args: args{
				in: &pb.UpdateUserDataRequest{
					Id:      1,
					Data:    []byte("test"),
					Version: 3,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "internal error",
			args: args{
				in: &pb.UpdateUserDataRequest{
					Id:      0,
					Data:    []byte("test"),
					Version: 3,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "no content",
			args: args{
				in: &pb.UpdateUserDataRequest{
					Id:      2,
					Data:    []byte("test"),
					Version: 3,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := srv
			got, err := s.UpdateUserData(ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateUserData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateUserData() got = %v, want %v", got, tt.want)
			}
		})
	}
}
