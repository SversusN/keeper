package grpcclient

import (
	"reflect"
	"testing"

	"github.com/SversusN/keeper/internal/client/models"
)

func TestClient_GetUserData(t *testing.T) {
	type args struct {
		model models.UserDataModel
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
				model: models.UserDataModel{
					ID: 1,
				},
			},
			want: &models.UserData{
				Name:     "testData",
				DataType: "password",
				Data:     []byte("test"),
				ID:       1,
				Version:  1,
			},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				model: models.UserDataModel{
					ID: 2,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testClient.GetUserData(tt.args.model)
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

func TestClient_GetUserDataList(t *testing.T) {
	tests := []struct {
		name    string
		want    []models.UserDataList
		wantErr bool
	}{
		{
			name: "success",
			want: []models.UserDataList{{
				ID:       1,
				Name:     "testData",
				DataType: "password",
				Version:  1,
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testClient.GetUserDataList()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserDataList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserDataList() got = %v, want %v", got, tt.want)
			}
		})
	}
}
