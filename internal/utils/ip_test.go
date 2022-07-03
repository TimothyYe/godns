package utils

import (
	"testing"

	"github.com/TimothyYe/godns/internal/settings"
)

func TestGetCurrentIP(t *testing.T) {
	conf := &settings.Settings{IPUrl: "https://myip.biturl.top"}
	ip, _ := GetCurrentIP(conf)

	if ip == "" {
		t.Log("IP is empty...")
	} else {
		t.Log("IP is:" + ip)
	}
}
