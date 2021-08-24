package ebmlpath

import "testing"

func TestMatch(t *testing.T) {
	type args struct {
		pattern string
		path    string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "match exact",
			args: args{pattern: `\a\global`, path: `\a\global`},
			want: true,
		},
		{
			name: "match placeholder 0-1 occurence",
			args: args{pattern: `\a\(0-1\)global`, path: `\a\global`},
			want: true,
		},
		{
			name: "match placeholder 0-1 occurence",
			args: args{pattern: `\a\(0-1\)global`, path: `\a\b\global`},
			want: true,
		},
		{
			name: "match placeholder 2+ occurence",
			args: args{pattern: `\a\(2-\)global`, path: `\a\b\c\d\global`},
			want: true,
		},
		{
			name: "match placeholder 1+ occurence",
			args: args{pattern: `\(1-\)global`, path: `\a\b\c\d\global`},
			want: true,
		},
		{
			name: "match recursive",
			args: args{pattern: `\a\+b\global`, path: `\a\b\global`},
			want: true,
		},
		{
			name: "match recursive",
			args: args{pattern: `\a\+b\global`, path: `\a\b\b\global`},
			want: true,
		},
		{
			name: "not match shorter pattern",
			args: args{pattern: `\a`, path: `\a\a`},
			want: false,
		},
		{
			name: "not match shorter path",
			args: args{pattern: `\a\a`, path: `\a`},
			want: false,
		},
		{
			name: "not match different",
			args: args{pattern: `\a\b\global`, path: `\a\global`},
			want: false,
		},
		{
			name: "not match placeholder 0-1 occurence",
			args: args{pattern: `\a\(0-1\)global`, path: `\a\b\c\global`},
			want: false,
		},
		{
			name: "not match placeholder 1-2 occurence",
			args: args{pattern: `\a\(1-2\)global`, path: `\a\global`},
			want: false,
		},
		{
			name: "not match placeholder 1+ occurence",
			args: args{pattern: `\a\(1-\)global`, path: `\a\global`},
			want: false,
		},
		{
			name: "not match placeholder 3+ occurence",
			args: args{pattern: `\a\(3-\)global`, path: `\a\global`},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Match(tt.args.pattern, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Match() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Match() got = %v, want %v", got, tt.want)
			}
		})
	}
}
