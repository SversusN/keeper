package app

import (
	"testing"

	"github.com/SversusN/keeper/internal/client/config"
	"github.com/SversusN/keeper/pkg/logger"
)

func TestClient_SaveData(t *testing.T) {
	type fields struct {
		gRPCClient   grpcClient
		printer      printer
		cache        clientCache
		Config       *config.Config
		Logger       *logger.Logger
		dataSyncChan chan int64
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "error",
			fields: fields{
				gRPCClient: &testGRPCClient{},
				printer:    &testPrinter{},
				cache:      &testCache{},
				Config:     testConfig,
				Logger:     testLogger,
			},
			wantErr: true,
		},
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
			if err := c.SaveData(); (err != nil) != tt.wantErr {
				t.Errorf("SaveData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_buildData(t *testing.T) {
	type args struct {
		dti int
		p   printer
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "build password",
			args: args{
				dti: 1,
				p:   &testPrinter{},
			},
			wantErr: false,
		},
		{
			name: "build card",
			args: args{
				dti: 2,
				p:   &testPrinter{},
			},
			wantErr: false,
		},
		{
			name: "build file",
			args: args{
				dti: 3,
				p:   &testPrinter{},
			},
			wantErr: true,
		},
		{
			name: "build text",
			args: args{
				dti: 4,
				p:   &testPrinter{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := buildData(tt.args.dti, tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
