package utils

import "testing"

func TestGetMD5Hash(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test1", args{"test"}, "098f6bcd4621d373cade4e832627b4f6"},
		{"test2", args{"test2"}, "ad0234829205b9033196ba818f7a872b"},
		{"test3", args{"test3"}, "8ad8757baa8564dc136c1e07507f4a98"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMD5Hash(tt.args.input); got != tt.want {
				t.Errorf("GetMD5Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}
