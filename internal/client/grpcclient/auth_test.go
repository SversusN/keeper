package grpcclient

import (
	"testing"

	"github.com/SversusN/keeper/internal/client/models"
)

func TestClient_SignIn_SignUp(t *testing.T) {
	type args struct {
		model models.AuthModel
	}
	tests := []struct {
		name    string
		args    args
		want    models.AuthToken
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				model: models.AuthModel{
					Login:    "testUser",
					Password: "testPass",
				},
			},
			want:    models.AuthToken("testToken"),
			wantErr: false,
		},
		{
			name: "unauthorized",
			args: args{
				model: models.AuthModel{
					Login:    "errorUser",
					Password: "testPass",
				},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testClient.SignIn(tt.args.model)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Register() got = %v, want %v", got, tt.want)
			}

			got2, err := testClient.Register(tt.args.model)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignUp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got2 != tt.want {
				t.Errorf("SignUp() got = %v, want %v", got2, tt.want)
			}
		})
	}
}
