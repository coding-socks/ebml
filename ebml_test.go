package ebml

import (
	"bytes"
	"io"
	"testing"
)

func TestReadElementID(t *testing.T) {
	type args struct {
		r           io.Reader
		maxIDLength uint
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "1 byte",
			args: args{r: bytes.NewReader([]byte{0x81}), maxIDLength: 8},
			want: "0x81",
		},
		{
			name: "2 byte",
			args: args{r: bytes.NewReader([]byte{0x41, 0x11}), maxIDLength: 8},
			want: "0x4111",
		},
		{
			name:    "early EOF",
			args:    args{r: bytes.NewReader([]byte{}), maxIDLength: 8},
			wantErr: true,
		},
		{
			name:    "early EOF",
			args:    args{r: bytes.NewReader([]byte{0x41}), maxIDLength: 8},
			wantErr: true,
		},
		{
			name:    "invalid length",
			args:    args{r: bytes.NewReader([]byte{0x41, 0x11}), maxIDLength: 1},
			wantErr: true,
		},
		{
			name:    "all one",
			args:    args{r: bytes.NewReader([]byte{0xff}), maxIDLength: 8},
			wantErr: true,
		},
		{
			name:    "all one",
			args:    args{r: bytes.NewReader([]byte{0x7f, 0xff}), maxIDLength: 8},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := ReadElementID(tt.args.r, tt.args.maxIDLength)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadElementID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadElementID() got = %v, want %v", got, tt.want)
			}
		})
	}
}
