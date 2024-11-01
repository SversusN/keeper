package handlers

import (
	"context"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	mock_storage "github.com/SversusN/keeper/internal/server/storage/mocks"
	pb "github.com/SversusN/keeper/pkg/grpc"
)

func TestServer_Ping(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB = mock_storage.NewMockrepository(ctrl)
	srv.Storage = mockDB
	mockDB.EXPECT().HealthCheck().Return(nil).AnyTimes()
	type args struct {
		in *pb.PingRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.PingResponse
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				in: &pb.PingRequest{},
			},
			want: &pb.PingResponse{
				Result: "OK",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := srv
			got, err := s.Ping(context.Background(), tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ping() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ping() got = %v, want %v", got, tt.want)
			}
		})
	}
}
