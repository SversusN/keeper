package grpcclient

import (
	"testing"

	"github.com/SversusN/keeper/internal/client/models"
)

func TestClient_UpdateUserData(t *testing.T) {
	type args struct {
		model *models.UserData
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				model: &models.UserData{
					Name:     "testData",
					DataType: "password",
					Data:     []byte("test"),
					Version:  1,
				},
			},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				model: &models.UserData{
					Name:     "testData",
					DataType: "password",
					Data:     []byte("test"),
					Version:  2,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testClient.UpdateUserData(tt.args.model); (err != nil) != tt.wantErr {
				t.Errorf("UpdateUserData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
