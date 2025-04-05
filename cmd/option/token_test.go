package option

import (
	"errors"
	"testing"
)

const (
	mocked_anonymous_registry = "anonymous_registry"
	mocked_auth_registry      = "auth_registry"
	mocked_invalid_registry   = "invalid_registry"
	mocked_identity_token     = "mocked_identity_token"
	mocked_registry_token     = "mocked_registry_token"
	mocked_invalid_token      = "invalid_token"
)

func TestParseTokenOption(t *testing.T) {
	// Mocking the getAuthToken function
	getAuthToken = func(registry string, identity_token string) (string, error) {
		switch {
		case registry == mocked_anonymous_registry:
			// mock returning an anonymous registry token
			return mocked_identity_token, nil
		case registry == mocked_auth_registry && identity_token == mocked_identity_token:
			// exchanged registry token
			return mocked_registry_token, nil
		}
		return "", errors.New("invalid registry")
	}

	type args struct {
		tokenOption string
		registry    string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Valid none option, no registry",
			args: args{
				tokenOption: "none",
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "Valid anonymous option, valid registry",
			args: args{
				tokenOption: "anonymous",
				registry:    mocked_anonymous_registry,
			},
			want:    mocked_identity_token,
			wantErr: false,
		},
		{
			name: "Valid anonymous option, invalid registry",
			args: args{
				tokenOption: "anonymous",
				registry:    mocked_invalid_registry,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Anonymous option with empty registry",
			args: args{
				tokenOption: "anonymous",
				registry:    "",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Valid exchange",
			args: args{
				tokenOption: "token=" + mocked_identity_token,
				registry:    mocked_auth_registry,
			},
			want:    mocked_registry_token,
			wantErr: false,
		},
		{
			name: "Invalid exchange",
			args: args{
				tokenOption: "token=" + mocked_invalid_token,
				registry:    mocked_auth_registry,
			},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Empty token option",
			args:    args{},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTokenOption(tt.args.tokenOption, tt.args.registry)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTokenOption() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseTokenOption() = %v, want %v", got, tt.want)
			}
		})
	}
}
