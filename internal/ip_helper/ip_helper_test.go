package ip_helper

import (
	"testing"

	"github.com/TimothyYe/godns/internal/settings"
)

func TestGetCurrentIP(t *testing.T) {
	conf := &settings.Settings{IPUrls: []string{"https://myip.biturl.top", "https://ip4.seeip.org"}}
	helper := NewIPHelper(conf)
	ip := helper.GetCurrentIP()

	if ip == "" {
		t.Log("IP is empty...")
	} else {
		t.Log("IP is:" + ip)
	}
}
