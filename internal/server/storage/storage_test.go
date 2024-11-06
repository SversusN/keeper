//go:build integration
// +build integration

package storage

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"

	"github.com/SversusN/keeper/pkg/logger"
)

var (
	testLogger   *logger.Logger
	testDB       *pgx.Conn
	testDatabase *Database
)

func TestMain(m *testing.M) {
	var err error
	testLogger, err = logger.Initialize("debug")
	if err != nil {
		panic(err)
	}

	//pool, err := dockertest.NewPool("")
	pool, err := dockertest.NewPool("unix:///Users/arseniy/.docker/run/docker.sock")
	if err != nil {
		testLogger.Log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		testLogger.Log.Fatalf("Could not connect to Docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "11",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=user_name",
			"POSTGRES_DB=dbname",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		testLogger.Log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseURL := fmt.Sprintf("postgres://user_name:secret@%s/dbname?sslmode=disable", hostAndPort)

	testLogger.Log.Info("Connecting to database on url: ", databaseURL)

	err = resource.Expire(60) // Tell docker to hard kill the container in 120 seconds
	if err != nil {
		testLogger.Log.Fatalf("Could not purge resource: %s", err)
	}

	// exponential backoff-retry, because the app in the container might not be ready to accept connections yet
	pool.MaxWait = 60 * time.Second
	if err = pool.Retry(func() error {
		testDB, err = pgx.Connect(context.Background(), databaseURL)
		if err != nil {
			return err
		}
		return testDB.Ping(context.Background())
	}); err != nil {
		testLogger.Log.Fatalf("Could not connect to docker: %s", err)
	}

	testDatabase = &Database{
		DB: testDB,
		l:  testLogger,
	}
	err = runMigrations(databaseURL)
	if err != nil {
		testLogger.Log.Fatalf("run migrations error: %s", err)
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		testLogger.Log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestDatabase_HealthCheck(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "without error",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testDatabase.HealthCheck(); (err != nil) != tt.wantErr {
				t.Errorf("HealthCheck() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabase_FindUserByLogin(t *testing.T) {
	type args struct {
		login string
	}
	tests := []struct {
		name        string
		args        args
		createLogin string
		want        *User
		wantErr     bool
	}{
		{
			name: "success",
			args: args{
				login: "testUser",
			},
			createLogin: "testUser",
			want: &User{
				Login:    "testUser",
				Password: "testPass",
				ID:       1,
			},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				login: "unknown",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.createLogin != "" {
				_, err := testDatabase.CreateUser(context.Background(), "testUser", "testPass")
				if err != nil {
					t.Errorf("CreateUser() error = %v", err)
				}
			}

			got, err := testDatabase.FindUserByLogin(context.Background(), tt.args.login)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindUserByLogin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindUserByLogin() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabase_CreateUser(t *testing.T) {
	type args struct {
		login    string
		password string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "already exists",
			args: args{
				login:    "testUser",
				password: "testPass",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testDatabase.CreateUser(context.Background(), tt.args.login, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CreateUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabase_SaveUserData(t *testing.T) {
	type args struct {
		userID   int64
		name     string
		dataType string
		data     []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				userID:   1,
				name:     "firstRecord",
				dataType: "password",
				data:     []byte("test"),
			},
			wantErr: false,
		},
		{
			name: "conflict",
			args: args{
				userID:   1,
				name:     "firstRecord",
				dataType: "password",
				data:     []byte("test"),
			},
			wantErr: true,
		},
		{
			name: "db error",
			args: args{
				userID:   1,
				name:     "firstRecord",
				dataType: "password",
				data:     nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testDatabase.SaveUserData(context.Background(), tt.args.userID, tt.args.name, tt.args.dataType, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("SaveUserData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabase_GetUserData(t *testing.T) {
	type args struct {
		userID int64
	}
	tests := []struct {
		name    string
		args    args
		want    []ShortRecord
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				userID: 1,
			},
			want: []ShortRecord{{
				Name:     "firstRecord",
				DataType: "password",
				ID:       1,
				Version:  1,
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testDatabase.GetUserData(context.Background(), tt.args.userID)
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

func TestDatabase_FindUserRecord(t *testing.T) {
	type args struct {
		id     int64
		userID int64
	}
	tests := []struct {
		name    string
		args    args
		want    *Record
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				id:     1,
				userID: 1,
			},
			want: &Record{
				Name:     "firstRecord",
				DataType: "password",
				Data:     []byte("test"),
				ID:       1,
				Version:  1,
			},
			wantErr: false,
		},
		{
			name: "no content",
			args: args{
				id:     1,
				userID: 5,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testDatabase.FindUserRecord(context.Background(), tt.args.id, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindUserRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindUserRecord() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabase_UpdateUserRecord(t *testing.T) {
	type args struct {
		rec *Record
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				rec: &Record{
					Name:     "firstRecord",
					DataType: "password",
					Data:     []byte("updated"),
					ID:       1,
					Version:  2,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testDatabase.UpdateUserRecord(context.Background(), tt.args.rec); (err != nil) != tt.wantErr {
				t.Errorf("UpdateUserRecord() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabase_Close(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "success",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testDatabase.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewStorage(t *testing.T) {
	type args struct {
		dsn string
		l   *logger.Logger
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "error",
			args: args{
				dsn: "",
				l:   testLogger,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewStorage(tt.args.dsn, tt.args.l)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewStorage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
