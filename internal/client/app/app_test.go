package app

import (
	"encoding/json"
	"net"
	"os"
	"testing"

	"github.com/SversusN/keeper/internal/client/cache"
	"github.com/SversusN/keeper/internal/client/config"
	"github.com/SversusN/keeper/internal/client/grpcclient"
	"github.com/SversusN/keeper/internal/client/models"
	"github.com/SversusN/keeper/pkg/logger"
)

var (
	testLogger *logger.Logger
	testConfig *config.Config
	testClient *Client
	testData   *models.UserData
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
		Config:       testConfig,
		Logger:       testLogger,
		gRPCClient:   &testGRPCClient{},
		printer:      &testPrinter{},
		cache:        &testCache{},
		dataSyncChan: make(chan int64),
	}

	passData, _ := json.Marshal(&PasswordData{
		Site:     "test.com",
		Login:    "testLogin",
		Password: "testPass",
	})
	testData = &models.UserData{
		Name:     "testName",
		DataType: "password",
		Data:     passData,
		ID:       1,
		Version:  1,
	}

	code := m.Run()

	os.Exit(code)
}

func TestNewClient(t *testing.T) {
	type args struct {
		l  *logger.Logger
		c  *config.Config
		bv string
		bd string
	}
	tests := []struct {
		name    string
		args    args
		grpcRun bool
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				l: &logger.Logger{},
				c: &config.Config{
					ChanSize:          1,
					ConnectionTimeout: 5,
					Host:              "localhost:3200",
				},
				bd: "",
				bv: "",
			},
			grpcRun: true,
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				l: &logger.Logger{},
				c: &config.Config{},
			},
			grpcRun: false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.grpcRun {
				conn, err := net.Listen("tcp", ":3200")
				if err != nil {
					t.Errorf("gRPC server start error = %v", err)
				}
				defer func() {
					_ = conn.Close()
				}()
			}

			_, err := NewClient(tt.args.l, tt.args.c, tt.args.bd, tt.args.bd)
			if (err != nil) && !tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

type testCache struct{}

type testGRPCClient struct{}

type testPrinter struct{}

func (t *testCache) Append(_ *models.UserData) {}

func (t *testCache) GetUserData(model models.UserDataModel) (*models.UserData, error) {
	if model.ID == 2 {
		return nil, cache.ErrNotFound
	}
	return testData, nil
}

func (t *testCache) GetUserDataList() []models.UserDataList {
	return []models.UserDataList{{
		Name:     "testName",
		DataType: "password",
		ID:       1,
		Version:  1,
	}}
}

func (t *testGRPCClient) Register(model models.AuthModel) (models.AuthToken, error) {
	if model.Login == "errorLogin" {
		return "", grpcclient.ErrRequest
	}
	return "testToken", nil
}

func (t *testGRPCClient) SignIn(model models.AuthModel) (models.AuthToken, error) {
	if model.Login == "errorLogin" {
		return "", grpcclient.ErrRequest
	}
	return "testToken", nil
}

func (t *testGRPCClient) GetUserData(model models.UserDataModel) (*models.UserData, error) {
	if model.ID == 2 {
		return nil, grpcclient.ErrRequest
	}
	return testData, nil
}

func (t *testGRPCClient) GetUserDataList() ([]models.UserDataList, error) {
	return []models.UserDataList{{
		Name:     "testName",
		DataType: "password",
		ID:       1,
		Version:  1,
	}}, nil
}

func (t *testGRPCClient) SaveUserData(model *models.UserData) error {
	if model.ID == 2 {
		return grpcclient.ErrRequest
	}
	return nil
}

func (t *testGRPCClient) UpdateUserData(model *models.UserData) error {
	if model.ID == 2 {
		return grpcclient.ErrRequest
	}
	return nil
}

func (t *testPrinter) Print(_ string) {}

func (t *testPrinter) Scan(a ...interface{}) (int, error) {
	return 1, nil
}
func (t *testPrinter) PrintLogo() {}
