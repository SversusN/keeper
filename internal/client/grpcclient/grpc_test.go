package grpcclient

import (
	"context"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/SversusN/keeper/internal/client/config"
	"github.com/SversusN/keeper/pkg/logger"

	pb "github.com/SversusN/keeper/pkg/grpc/keeper"
)

var (
	testLogger *logger.Logger
	testConfig *config.Config
	testClient *Client
)

func TestMain(m *testing.M) {
	var err error
	testLogger, err = logger.Initialize("info")
	if err != nil {
		panic(err)
	}

	testConfig = &config.Config{
		Host:              "localhost:3200",
		LogLevel:          "debug",
		ConnectionTimeout: 1,
		ChanSize:          1,
	}

	testClient = &Client{
		config:     testConfig,
		gRPCClient: &testGRPCClient{},
		timeout:    5 * time.Second,
	}

	code := m.Run()

	os.Exit(code)
}

func TestClient_GetAuthToken(t *testing.T) {
	type fields struct {
		authToken string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "return authToken",
			fields: fields{authToken: "testToken"},
			want:   "testToken",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				authToken: tt.fields.authToken,
			}
			if got := c.GetAuthToken(); got != tt.want {
				t.Errorf("GetAuthToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewGRPCClient(t *testing.T) {
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		testLogger.Log.Error(err)
		t.Errorf("tcp connection failed: %v", err)
	}
	defer func() {
		_ = listen.Close()
	}()

	type args struct {
		c *config.Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				c: testConfig,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewGRPCClient(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGRPCClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

type testGRPCClient struct{}

func (t *testGRPCClient) Ping(ctx context.Context, in *pb.PingRequest, opts ...grpc.CallOption) (
	*pb.PingResponse, error) {
	return &pb.PingResponse{}, nil
}

func (t *testGRPCClient) SignUp(ctx context.Context, in *pb.SignUpRequest, opts ...grpc.CallOption) (
	*pb.SignUpResponse, error) {
	if in.Login == "errorUser" {
		return nil, status.Error(codes.Unauthenticated, http.StatusText(http.StatusUnauthorized))
	}
	return &pb.SignUpResponse{
		Token: "testToken",
	}, nil
}

func (t *testGRPCClient) SignIn(ctx context.Context, in *pb.SignInRequest, opts ...grpc.CallOption) (
	*pb.SignInResponse, error) {
	if in.Login == "errorUser" {
		return nil, status.Error(codes.Unauthenticated, http.StatusText(http.StatusUnauthorized))
	}
	return &pb.SignInResponse{
		Token: "testToken",
	}, nil
}

func (t *testGRPCClient) SaveData(ctx context.Context, in *pb.SaveDataRequest, opts ...grpc.CallOption) (
	*pb.SaveDataResponse, error) {
	if in.DataType == "card" {
		return nil, status.Error(codes.Internal, http.StatusText(http.StatusInternalServerError))
	}
	return &pb.SaveDataResponse{
		Result: "OK",
	}, nil
}

func (t *testGRPCClient) GetUserDataList(ctx context.Context, in *pb.UserDataListRequest, opts ...grpc.CallOption) (
	*pb.UserDataListResponse, error) {
	return &pb.UserDataListResponse{
		Data: []*pb.UserDataNested{{
			Id:       1,
			Name:     "testData",
			DataType: "password",
			Version:  1,
		}},
	}, nil
}

func (t *testGRPCClient) GetUserData(ctx context.Context, in *pb.UserDataRequest, opts ...grpc.CallOption) (
	*pb.UserDataResponse, error) {
	if in.Id == 2 {
		return nil, status.Error(codes.Internal, http.StatusText(http.StatusInternalServerError))
	}
	return &pb.UserDataResponse{
		Name:     "testData",
		DataType: "password",
		Data:     []byte("test"),
		Id:       1,
		Version:  1,
	}, nil
}

func (t *testGRPCClient) UpdateUserData(ctx context.Context, in *pb.UpdateUserDataRequest, opts ...grpc.CallOption) (
	*pb.UpdateUserDataResponse, error) {
	if in.Version == 2 {
		return nil, status.Error(codes.Internal, http.StatusText(http.StatusInternalServerError))
	}
	return &pb.UpdateUserDataResponse{
		Result: "OK",
	}, nil
}
