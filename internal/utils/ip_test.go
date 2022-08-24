package utils

import (
	"testing"

	"github.com/TimothyYe/godns/internal/settings"
)

func TestGetCurrentIP(t *testing.T) {
	conf := &settings.Settings{IPUrls: []string{"https://aaa.bbb.ccc", "https://myip.biturl.top", "https://ip4.seeip.org"}}
	ip, _ := GetCurrentIP(conf)

	if ip == "" {
		t.Log("IP is empty...")
	} else {
		t.Log("IP is:" + ip)
	}
}
