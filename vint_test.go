package ebml

import "testing"

func Test_vintOctetLength(t *testing.T) {
	tests := []struct {
		name string
		vint []byte
		want int
	}{
		{
			name: "1", vint: []byte{0x80}, want: 1,
		},
		{
			name: "2", vint: []byte{0x40}, want: 2,
		},
		{
			name: "7", vint: []byte{0x02}, want: 7,
		},
		{
			name: "8", vint: []byte{0x01}, want: 8,
		},
		{
			name: "9", vint: []byte{0x00, 0x80}, want: 9,
		},
		{
			name: "10", vint: []byte{0x00, 0x40}, want: 10,
		},
		{
			name: "15", vint: []byte{0x00, 0x02}, want: 15,
		},
		{
			name: "16", vint: []byte{0x00, 0x01}, want: 16,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := vintOctetLength(tt.vint); got != tt.want {
				t.Errorf("vintOctetLength() = %v, want %v", got, tt.want)
			}
		})
	}
}
