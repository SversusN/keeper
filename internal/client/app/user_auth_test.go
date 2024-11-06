package app

import (
	"reflect"
	"testing"

	"github.com/SversusN/keeper/internal/client/models"
)

func TestClient_userSignUp_userSignIn(t *testing.T) {
	type args struct {
		authM models.AuthModel
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				authM: models.AuthModel{
					Login:    "testLogin",
					Password: "testPass",
				},
			},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				authM: models.AuthModel{
					Login:    "badLogin",
					Password: "badPass",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testClient.register(tt.args.authM); (err != nil) && tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := testClient.signIn(tt.args.authM); (err != nil) && tt.wantErr {
				t.Errorf("signIn() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_buildAuthData(t *testing.T) {
	type args struct {
		p printer
	}
	tests := []struct {
		name    string
		args    args
		want    *models.AuthModel
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				p: &testPrinter{},
			},
			want: &models.AuthModel{
				Login:    "",
				Password: "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildAuthData(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildAuthData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildAuthData() got = %v, want %v", got, tt.want)
			}
		})
	}
}
