package handlers

import (
	"context"
	"fmt"
	"github.com/SversusN/keeper/internal/utils"
	"github.com/SversusN/keeper/internal/utils/encrypter"
	"os"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/SversusN/keeper/internal/server/config"
	"github.com/SversusN/keeper/internal/server/storage"
	mock_storage "github.com/SversusN/keeper/internal/server/storage/mocks"
	"github.com/SversusN/keeper/pkg/logger"
)

var (
	testLogger *logger.Logger
	testConfig *config.Config
	testUser   *storage.User
	testRecord *storage.Record
	ctx        = context.WithValue(context.Background(), utils.UserIDContextKey, int64(1))
	srv        *Server
	mockDB     *mock_storage.Mockrepository
)

func TestMain(m *testing.M) {
	var err error
	testLogger, err = logger.Initialize("info")
	if err != nil {
		panic(err)
	}

	testConfig = &config.Config{
		DatabaseDSN: "",
		Host:        "localhost:3200",
		LogLevel:    "debug",
		SecretKey:   "test",
	}

	srv = &Server{
		Storage: mockDB,
		crypto:  &testCrypt{},
		Config:  testConfig,
		Logger:  testLogger,
	}

	testUser = &storage.User{
		Login:    "testUser",
		Password: "testPassword",
		ID:       0,
	}
	testRecord = &storage.Record{
		Name:     "testName",
		DataType: "password",
		Data:     []byte("test"),
		ID:       1,
		Version:  1,
	}

	code := m.Run()

	os.Exit(code)
}

func TestNewServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB = mock_storage.NewMockrepository(ctrl)
	srv.Storage = mockDB

	type args struct {
		r repository
		c *config.Config
		l *logger.Logger
	}
	tests := []struct {
		name string
		args args
		want *Server
	}{
		{
			name: "success",
			args: args{
				r: mockDB,
				c: testConfig,
				l: testLogger,
			},
			want: &Server{
				Storage: mockDB,
				crypto:  &encrypter.Token{},
				Config:  testConfig,
				Logger:  testLogger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewServer(tt.args.r, tt.args.c, tt.args.l); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

type testCrypt struct{}

func (t *testCrypt) HashFunc(src string) (string, error) {
	return "test_hash", nil
}

func (t *testCrypt) CompareHash(src, hash string) error {
	if src == "errPass" {
		return fmt.Errorf("not equal")
	}
	return nil
}

func (t *testCrypt) BuildJWT(userID int64, secret string) (string, error) {
	return "test_token", nil
}

func (t *testCrypt) GetUserID(tokenString, secret string) (int64, error) {
	return 1, nil
}
