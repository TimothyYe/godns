package lib

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/TimothyYe/godns/internal/settings"
	"github.com/TimothyYe/godns/internal/utils"
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

// TestGetIPOnlineRejectsWrongFamily checks that getIPOnline only hands back an
// address whose family matches the configured ip_type. When every attempt is
// skipped the result must be empty so getCurrentIP can fall back to the
// interface lookup instead of caching an address of the wrong family.
func TestGetIPOnlineRejectsWrongFamily(t *testing.T) {
	cases := []struct {
		name   string
		ipType string
		answer string
		want   string
	}{
		{name: "IPv6 mode skips an IPv4 answer", ipType: utils.IPV6, answer: "1.2.3.4\n", want: ""},
		{name: "IPv4 mode skips an IPv6 answer", ipType: utils.IPV4, answer: "2001:db8::1\n", want: ""},
		{name: "IPv6 mode keeps an IPv6 answer", ipType: utils.IPV6, answer: "2001:db8::1\n", want: "2001:db8::1"},
		{name: "IPv4 mode keeps an IPv4 answer", ipType: utils.IPV4, answer: "1.2.3.4\n", want: "1.2.3.4"},
		{name: "unknown ip_type skips every answer", ipType: "IPv4/IPv6", answer: "2001:db8::1\n", want: ""},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				fmt.Fprint(w, c.answer)
			}))
			defer server.Close()

			helper := &IPHelper{
				configuration: &settings.Settings{IPType: c.ipType},
				reqURLs:       []string{server.URL},
				idx:           -1,
			}

			got := helper.getIPOnline()
			if got != c.want {
				t.Fatalf("ip_type %q with endpoint answer %q: getIPOnline() = %q, want %q",
					c.ipType, strings.TrimSpace(c.answer), got, c.want)
			}
		})
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
