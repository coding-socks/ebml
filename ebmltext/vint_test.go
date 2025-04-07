package ebmltext

import (
	"bytes"
	"testing"
)

func Test_vintDataAllOne(t *testing.T) {
	type args struct {
		b uint64
		l int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "1 byte all 1",
			args: args{b: 0x7f, l: 1},
			want: true,
		},
		{
			name: "1 byte not all 1",
			args: args{b: 0x70, l: 1},
			want: false,
		},
		{
			name: "2 byte all 1",
			args: args{b: 0x3fff, l: 2},
			want: true,
		},
		{
			name: "2 byte not all 1",
			args: args{b: 0x3ff0, l: 2},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := vintDataAllOne(tt.args.b, tt.args.l); got != tt.want {
				t.Errorf("vintDataAllOne() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_vint(t *testing.T) {
	type args struct {
		vint uint64
		l    int
	}
	tests := []struct {
		name     string
		args     args
		wantBuf  []byte
		wantVint uint64
	}{
		{
			name:     "1 byte",
			args:     args{vint: 0x02, l: 1},
			wantBuf:  []byte{0x82},
			wantVint: 0x02,
		},
		{
			name:     "1 byte on 2 bytes",
			args:     args{vint: 0x02, l: 2},
			wantBuf:  []byte{0x40, 0x02},
			wantVint: 0x02,
		},
		{
			name:     "2 bytes",
			args:     args{vint: 0x0111, l: 2},
			wantBuf:  []byte{0x41, 0x11},
			wantVint: 0x0111,
		},
		{
			name:     "7 bytes",
			args:     args{vint: 0x01111111111111, l: 7},
			wantBuf:  []byte{0x03, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11},
			wantVint: 0x01111111111111,
		},
		{
			name:     "8 bytes",
			args:     args{vint: 0x0011111111111111, l: 8},
			wantBuf:  []byte{0x01, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11},
			wantVint: 0x0011111111111111,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := make([]byte, 8)
			l, _ := AppendVintData(tt.args.vint, tt.args.l, buf)
			if got := buf[:l]; !bytes.Equal(buf[:l], tt.wantBuf) {
				t.Errorf("AppendVintData() = %x, want %x", got[:l], tt.wantBuf)
			}
			if got, _, _ := ReadVintData(buf); got != tt.wantVint {
				t.Errorf("ReadVintData() = %v, want %v", got, tt.wantVint)
			}
		})
	}
}
