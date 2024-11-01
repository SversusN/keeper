package grpcclient

import (
	"testing"

	"github.com/SversusN/keeper/internal/client/models"
)

func TestClient_SaveUserData(t *testing.T) {
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
					Name:     "testName",
					DataType: "password",
					Data:     []byte("test"),
				},
			},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				model: &models.UserData{
					Name:     "testName",
					DataType: "card",
					Data:     []byte("test"),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testClient.SaveUserData(tt.args.model); (err != nil) != tt.wantErr {
				t.Errorf("SaveUserData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
