package option

import (
	"errors"
	"testing"
)

const (
	anonymous_registry = "anonymous_registry"
	mocked_token       = "mocked_token"
	invalid_registry   = "invalid_registry"
)

func TestParseTokenOption(t *testing.T) {
	getAuthToken = func(registry string) (string, error) {
		if registry == anonymous_registry {
			return mocked_token, nil
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
				registry:    anonymous_registry,
			},
			want:    mocked_token,
			wantErr: false,
		},
		{
			name: "Valid anonymous option, invalid registry",
			args: args{
				tokenOption: "anonymous",
				registry:    invalid_registry,
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
			name: "Valid token option",
			args: args{
				tokenOption: "token=mytoken",
				registry:    "example.com",
			},
			want:    "mytoken",
			wantErr: false,
		},
		{
			name: "Invalid token option",
			args: args{
				tokenOption: "invalid",
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
