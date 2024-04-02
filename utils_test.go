package provider

import (
	"testing"
)

func Test_basePath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{{
		name: "main",
		args: args{
			path: "https://yourservice.com/clevercloud/resources",
		},
		want: "/clevercloud/resources",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := basePath(tt.args.path); got != tt.want {
				t.Errorf("basePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
