package app

import (
	"reflect"
	"testing"

	"github.com/SversusN/keeper/internal/client/config"
	"github.com/SversusN/keeper/internal/client/models"
	"github.com/SversusN/keeper/pkg/logger"
)

func TestClient_GetDataFromCache(t *testing.T) {
	type args struct {
		m models.UserDataModel
	}
	tests := []struct {
		name    string
		args    args
		want    *models.UserData
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				m: models.UserDataModel{ID: 1},
			},
			want:    testData,
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				m: models.UserDataModel{ID: 2},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testClient.GetDataFromCache(tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDataFromCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDataFromCache() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetDataListFromCache(t *testing.T) {
	tests := []struct {
		name string
		want []models.UserDataList
	}{
		{
			name: "success",
			want: []models.UserDataList{{
				Name:     "testName",
				DataType: "password",
				ID:       1,
				Version:  1,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := testClient.GetDataListFromCache(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDataListFromCache() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_SyncCache(t *testing.T) {
	type fields struct {
		gRPCClient   grpcClient
		printer      printer
		cache        clientCache
		Config       *config.Config
		Logger       *logger.Logger
		dataSyncChan chan int64
	}
	type args struct {
		records []models.UserDataList
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				gRPCClient:   tt.fields.gRPCClient,
				printer:      tt.fields.printer,
				cache:        tt.fields.cache,
				Config:       tt.fields.Config,
				Logger:       tt.fields.Logger,
				dataSyncChan: tt.fields.dataSyncChan,
			}
			c.SyncCache(tt.args.records)
		})
	}
}

func TestClient_UpdateDataInCache(t *testing.T) {
	type args struct {
		data *models.UserData
	}
	tests := []struct {
		name     string
		args     args
		wantSize int
	}{
		{
			name: "success",
			args: args{
				data: testData,
			},
			wantSize: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testClient.UpdateDataInCache(tt.args.data)
			size := len(testClient.cache.GetUserDataList())
			if size != tt.wantSize {
				t.Errorf("UpdateDataInCache error, current size: %v", size)
			}
		})
	}
}
