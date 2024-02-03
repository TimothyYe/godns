package lib

import (
	"testing"

	"github.com/TimothyYe/godns/internal/settings"
	"github.com/TimothyYe/godns/pkg/lib"
)

func TestGetCurrentIP(t *testing.T) {
	t.Skip()
	conf := &settings.Settings{IPUrls: []string{"https://myip.biturl.top"}}
	helper := lib.GetIPHelperInstance(conf)
	ip := helper.GetCurrentIP()

	if ip == "" {
		t.Log("IP is empty...")
	} else {
		t.Log("IP is:" + ip)
	}
}
