package app

import "testing"

func TestClient_EditData(t *testing.T) {
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
			if err := testClient.EditData(); (err != nil) != tt.wantErr {
				t.Errorf("EditData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
