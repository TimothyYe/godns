package godns

import (
	"testing"
)

func testGetCurrentIP(t *testing.T) {
	conf := &Settings{IPUrl: "http://members.3322.org/dyndns/getip"}
	ip, _ := GetCurrentIP(conf)

	if ip == "" {
		t.Log("IP is empty...")
	} else {
		t.Log("IP is:" + ip)
	}
}
