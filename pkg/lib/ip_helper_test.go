package lib

import (
	"testing"
	"time"

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

func TestGetMikrotikIP(t *testing.T) {
	t.Skip()

	conf := &settings.Settings{
		Mikrotik: settings.Mikrotik{
			Enabled:   true,
			Addr:      "http://192.168.20.1:81",
			Username:  "admin",
			Password:  "",
			Interface: "pppoe-out",
		},
	}
	helper := GetIPHelperInstance(conf)
	ip := helper.GetCurrentIP()

	if ip == "" {
		t.Log("IP is empty...")
	} else {
		t.Log("IP is:" + ip)
	}
}

// TestIPHelperStop verifies that calling Stop() closes the helper's
// internal stop channel (so the background refresh goroutine exits) and
// that Stop is safely idempotent. Note: this test mutates the package
// singleton and must be the last live test to touch the helper —
// after Stop the singleton's refresh goroutine cannot be restarted.
func TestIPHelperStop(t *testing.T) {
	helper := GetIPHelperInstance(&settings.Settings{Interval: 60})

	helper.Stop()

	select {
	case <-helper.stopCh:
		// closed as expected
	case <-time.After(time.Second):
		t.Fatal("stopCh was not closed after Stop()")
	}

	// Second Stop must be safe — sync.Once guards the close.
	helper.Stop()
}
