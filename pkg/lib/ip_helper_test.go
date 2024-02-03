package lib

import (
	"testing"

	"github.com/TimothyYe/godns/internal/settings"
)

func TestGetCurrentIP(t *testing.T) {
	t.Skip()
	conf := &settings.Settings{IPUrls: []string{"https://myip.biturl.top"}}
	helper := GetIPHelperInstance(conf)
	ip := helper.GetCurrentIP()

	if ip == "" {
		t.Log("IP is empty...")
	} else {
		t.Log("IP is:" + ip)
	}
}
