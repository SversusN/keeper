package handlers

import (
	"context"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/SversusN/keeper/internal/server/storage"
	mock_storage "github.com/SversusN/keeper/internal/server/storage/mocks"
	pb "github.com/SversusN/keeper/pkg/grpc"
)

func TestServer_SignIn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB = mock_storage.NewMockrepository(ctrl)
	srv.Storage = mockDB
	mockDB.EXPECT().FindUserByLogin(gomock.Any(), "noUser").Return(nil, storage.ErrNowRows).AnyTimes()
	mockDB.EXPECT().FindUserByLogin(gomock.Any(), "errUser").Return(nil, storage.ErrFindUser).AnyTimes()
	mockDB.EXPECT().FindUserByLogin(gomock.Any(), gomock.Any()).Return(testUser, nil).AnyTimes()

	type args struct {
		in *pb.SignInRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.SignInResponse
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				in: &pb.SignInRequest{
					Login:    "testUser",
					Password: "testPassword",
				},
			},
			want:    &pb.SignInResponse{Token: "test_token"},
			wantErr: false,
		},
		{
			name: "missing user",
			args: args{
				in: &pb.SignInRequest{
					Login:    "noUser",
					Password: "testPassword",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "unauthorized",
			args: args{
				in: &pb.SignInRequest{
					Login:    "testUser",
					Password: "errPass",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "internal error",
			args: args{
				in: &pb.SignInRequest{
					Login:    "errUser",
					Password: "testPassword",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := srv
			got, err := s.SignIn(context.Background(), tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Register() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_SignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB = mock_storage.NewMockrepository(ctrl)
	srv.Storage = mockDB
	mockDB.EXPECT().CreateUser(gomock.Any(), "conflictUser", gomock.Any()).Return(int64(0), storage.ErrConflict).AnyTimes()
	mockDB.EXPECT().CreateUser(gomock.Any(), "errUser", gomock.Any()).Return(int64(0), storage.ErrCreateUser).AnyTimes()
	mockDB.EXPECT().CreateUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(1), nil).AnyTimes()

	type args struct {
		in *pb.RegisterRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.RegisterResponse
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				in: &pb.RegisterRequest{
					Login:    "testUser",
					Password: "testPassword",
				},
			},
			want:    &pb.RegisterResponse{Token: "test_token"},
			wantErr: false,
		},
		{
			name: "conflict",
			args: args{
				in: &pb.RegisterRequest{
					Login:    "conflictUser",
					Password: "testPassword",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "internal error",
			args: args{
				in: &pb.RegisterRequest{
					Login:    "errUser",
					Password: "testPassword",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := srv
			got, err := s.SignUp(context.Background(), tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Register() got = %v, want %v", got, tt.want)
			}
		})
	}
}
