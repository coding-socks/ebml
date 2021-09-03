package ebml

import (
	"reflect"
	"testing"
)

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

func Test_vintData(t *testing.T) {
	type args struct {
		b []byte
		l int
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "1 byte",
			args: args{b: []byte{0x81}, l: 1},
			want: []byte{0x01},
		},
		{
			name: "2 byte",
			args: args{b: []byte{0x41, 0x11}, l: 2},
			want: []byte{0x01, 0x11},
		},
		{
			name: "7 byte",
			args: args{b: []byte{0x03, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}, l: 7},
			want: []byte{0x01, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11},
		},
		{
			name: "8 byte",
			args: args{b: []byte{0x01, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}, l: 8},
			want: []byte{0x00, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11},
		},
		{
			name: "9 byte",
			args: args{b: []byte{0x00, 0x81, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}, l: 9},
			want: []byte{0x00, 0x01, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11},
		},
		{
			name: "10 byte",
			args: args{b: []byte{0x00, 0x41, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}, l: 10},
			want: []byte{0x00, 0x01, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11},
		},
		{
			name: "15 byte",
			args: args{b: []byte{0x00, 0x03, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}, l: 15},
			want: []byte{0x00, 0x01, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11},
		},
		{
			name: "16 byte",
			args: args{b: []byte{0x00, 0x01, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}, l: 16},
			want: []byte{0x00, 0x00, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := vintData(tt.args.b, tt.args.l); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("vintData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_vintDataAllOne(t *testing.T) {
	type args struct {
		b []byte
		l int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "1 byte all 1",
			args: args{b: []byte{0x7f}, l: 1},
			want: true,
		},
		{
			name: "1 byte not all 1",
			args: args{b: []byte{0x70}, l: 1},
			want: false,
		},
		{
			name: "2 byte all 1",
			args: args{b: []byte{0x3f, 0xff}, l: 2},
			want: true,
		},
		{
			name: "2 byte not all 1",
			args: args{b: []byte{0x3f, 0xf0}, l: 2},
			want: false,
		},
		{
			name: "9 byte all 1",
			args: args{b: []byte{0x00, 0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, l: 9},
			want: true,
		},
		{
			name: "9 byte not all 1",
			args: args{b: []byte{0x00, 0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xf0}, l: 9},
			want: false,
		},
		{
			name: "10 byte all 1",
			args: args{b: []byte{0x00, 0x3f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, l: 10},
			want: true,
		},
		{
			name: "10 byte not all 1",
			args: args{b: []byte{0x00, 0x3f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xf0}, l: 10},
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
