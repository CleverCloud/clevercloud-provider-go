package provider

import (
	"testing"
	"time"
)

func Test_tokenSignature(t *testing.T) {
	type args struct {
		addonID   string
		salt      string
		timestamp string
	}
	tests := []struct {
		name string
		args args
		want string
	}{{
		name: "main",
		args: args{
			addonID:   "dummy_sd675fa67sf57asdf67asf57",
			salt:      "6hjj896df7g5sg6d5g9df5gsdg67",
			timestamp: "1712063189",
		},
		want: "1e47e0d5c1095fed750acef9411a23fcb85de43f",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//t.Errorf("%s%s%s", tt.args.addonID, tt.args.salt, tt.args.timestamp)
			if got := tokenSignature(tt.args.addonID, tt.args.salt, tt.args.timestamp); got != tt.want {
				t.Errorf("tokenSignature() = %s, want %s", got, tt.want)
			}
		})
	}
}

func Test_isOutdated(t *testing.T) {
	tests := []struct {
		name string
		t    int64
		want bool
	}{{
		name: "main", t: time.Now().Add(-5 * time.Minute).Unix(), want: false,
	}, {
		name: "outdated", t: time.Now().Add(-20 * time.Minute).Unix(), want: true,
	}}

	t.Logf("TEST %+v", time.UnixMilli(1712065143000).String())
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isOutdated(tt.t); got != tt.want {
				t.Errorf("isOutdated() = %v, want %v", got, tt.want)
			}
		})
	}
}
