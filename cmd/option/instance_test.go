package option

import (
	"testing"
	"time"
)

func TestParseInstanceOption(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Instance
		wantErr bool
	}{
		{
			name:  "Valid input with batch size and interval",
			input: "10=5/2s",
			want: Instance{
				Count:         10,
				BatchSize:     5,
				BatchInterval: 2 * time.Second,
			},
			wantErr: false,
		},
		{
			name:  "Valid input with single instance",
			input: "1=1/1s",
			want: Instance{
				Count:         1,
				BatchSize:     1,
				BatchInterval: 1 * time.Second,
			},
			wantErr: false,
		},

		{
			name:  "Valid input: no batch size and interval",
			input: "10",
			want: Instance{
				Count: 10,
			},
			wantErr: false,
		},
		{
			name:    "Invalid input: missing batch size and interval",
			input:   "10=",
			wantErr: true,
		},
		{
			name:    "Invalid input: missing batch size",
			input:   "10=/2s",
			wantErr: true,
		},
		{
			name:    "Invalid input: missing interval",
			input:   "10=2/",
			wantErr: true,
		},
		{
			name:    "Invalid input: non-numeric count",
			input:   "abc=5/2s",
			wantErr: true,
		},
		{
			name:    "Invalid input: non-numeric batch size",
			input:   "10=abc/2s",
			wantErr: true,
		},
		{
			name:    "Invalid input: invalid interval format",
			input:   "10=5/abc",
			wantErr: true,
		},
		{
			name:    "Invalid input: count is zero",
			input:   "0=5/2s",
			wantErr: true,
		},
		{
			name:    "Invalid input: batch size is zero",
			input:   "10=0/2s",
			wantErr: true,
		},
		{
			name:    "Invalid input: interval is zero",
			input:   "10=5/0s",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instance := &Instance{}
			instance.SetFlag(tt.input)
			err := instance.Parse()
			if tt.wantErr {
				if err == nil {
					t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else {
				if instance.Count != tt.want.Count {
					t.Errorf("Parse() Count = %v, want %v", instance.Count, tt.want.Count)
				}
				if instance.BatchSize != tt.want.BatchSize {
					t.Errorf("Parse() BatchSize = %v, want %v", instance.BatchSize, tt.want.BatchSize)
				}
				if instance.BatchInterval != tt.want.BatchInterval {
					t.Errorf("Parse() BatchInterval = %v, want %v", instance.BatchInterval, tt.want.BatchInterval)
				}
			}
		})
	}
}
