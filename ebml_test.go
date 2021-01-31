package ebml

import (
	"github.com/coding-socks/ebml/vint"
	"testing"
)

func Test_validateID(t *testing.T) {
	type args struct {
		id *vint.Vint
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// Source https://www.rfc-editor.org/rfc/rfc8794#name-element-id
		{
			name:    "all zero",
			args:    args{id: vint.NewVint([]byte{0b10000000})},
			wantErr: true,
		},
		{
			name:    "all zero",
			args:    args{id: vint.NewVint([]byte{0b01000000, 0b00000000})},
			wantErr: true,
		},
		{
			name:    "not all zero",
			args:    args{id: vint.NewVint([]byte{0b10000001})},
			wantErr: false,
		},
		{
			name:    "shorter available",
			args:    args{id: vint.NewVint([]byte{0b01000000, 0b00000001})},
			wantErr: true,
		},
		{
			name:    "shorter not available",
			args:    args{id: vint.NewVint([]byte{0b10111111})},
			wantErr: false,
		},
		{
			name:    "shorter available",
			args:    args{id: vint.NewVint([]byte{0b01000000, 0b00111111})},
			wantErr: true,
		},
		{
			name:    "all one",
			args:    args{id: vint.NewVint([]byte{0b11111111})},
			wantErr: true,
		},
		{
			name:    "not all one",
			args:    args{id: vint.NewVint([]byte{0b01000000, 0b01111111})},
			wantErr: false,
		},
		// Extra
		{
			name:    "all one",
			args:    args{id: vint.NewVint([]byte{0b01111111, 0b11111111})},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateID(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("validateID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
