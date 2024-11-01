package encrypter

import (
	"testing"
)

func TestToken_BuildJWT(t *testing.T) {
	type args struct {
		userID int64
		secret string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				userID: 1,
				secret: "test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Token{}
			_, err := c.BuildJWT(tt.args.userID, tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestToken_CompareHash(t *testing.T) {
	type args struct {
		src  string
		hash string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				src:  "testPass",
				hash: "$2a$10$E9Bl/QO1eb/M8ka0veg7IuZZbGtov/85.1VB2CVuk8YAes3e8x0tS",
			},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				src:  "errorPass",
				hash: "$2a$10$E9Bl/QO1eb/M8ka0veg7IuZZbGtov/85.1VB2CVuk8YAes3e8x0tS",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Token{}
			if err := c.CompareHash(tt.args.src, tt.args.hash); (err != nil) != tt.wantErr {
				t.Errorf("CompareHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestToken_GetUserID(t *testing.T) {
	type args struct {
		tokenString string
		secret      string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOjF9.-LWjc9k-f3WQmNzz43BQTPcFwCwjjrqfy2MpsAddeCw",
				secret:      "test",
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				tokenString: "fault",
				secret:      "test",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Token{}
			got, err := c.GetUserID(tt.args.tokenString, tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetUserID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncrypter_HashFunc(t *testing.T) {
	type args struct {
		src string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "success",
			args:    args{src: "testPass"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Token{}
			_, err := c.HashFunc(tt.args.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashFunc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
