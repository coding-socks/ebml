package ebml

import (
	"testing"
)

func Test_validateIDData(t *testing.T) {
	tests := []struct {
		name    string
		id      []byte
		wantErr bool
	}{
		// Source https://www.rfc-editor.org/rfc/rfc8794#name-element-id
		{
			name:    "all zero",
			id:      []byte{0b10000000},
			wantErr: false,
		},
		{
			name:    "all zero",
			id:      []byte{0b01000000, 0b00000000},
			wantErr: false,
		},
		{
			name:    "not all zero",
			id:      []byte{0b10000001},
			wantErr: false,
		},
		{
			name:    "all one",
			id:      []byte{0b11111111},
			wantErr: true,
		},
		{
			name:    "not all one",
			id:      []byte{0b01000000, 0b01111111},
			wantErr: false,
		},
		// Extra
		{
			name:    "all one",
			id:      []byte{0b01111111, 0b11111111},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := vintOctetLength(tt.id)
			if err := validateIDData(vintData(tt.id, w), w); (err != nil) != tt.wantErr {
				t.Errorf("validateIDData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
