package auth

import (
	"testing"
)

func TestPasswordPayload_Validate(t *testing.T) {
	tests := []struct {
		name    string
		payload PasswordPayload
		wantErr bool
	}{
		{
			name:    "Matching passwords",
			payload: PasswordPayload{Password: "strongpass", ConfirmPassword: "strongpass"},
			wantErr: false,
		},
		{
			name:    "Mismatched passwords",
			payload: PasswordPayload{Password: "strongpass", ConfirmPassword: "wrongpass"},
			wantErr: true,
		},
		{
			name:    "Empty passwords",
			payload: PasswordPayload{Password: "", ConfirmPassword: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.payload.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
