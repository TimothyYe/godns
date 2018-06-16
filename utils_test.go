package godns

import (
	"os"
	"strings"
	"testing"
)

func TestGetCurrentIP(t *testing.T) {
	conf := &Settings{IPUrl: "http://members.3322.org/dyndns/getip"}
	ip, _ := GetCurrentIP(conf)

	if ip == "" {
		t.Log("IP is empty...")
	} else {
		t.Log("IP is:" + ip)
	}

	conf = &Settings{Socks5Proxy: "localhost:8899", IPUrl: "http://members.3322.org/dyndns/getip"}
	ip, err := GetCurrentIP(conf)

	if ip != "" && err == nil {
		t.Error("should return error")
	}
}

func TestCheckSettings(t *testing.T) {
	settingError := &Settings{}
	if err := CheckSettings(settingError); err == nil {
		t.Error("setting is invalid, should return error")
	}

	settingDNSPod := &Settings{Provider: "DNSPod", LoginToken: "aaa"}
	if err := CheckSettings(settingDNSPod); err == nil {
		t.Log("setting with login token, passed")
	} else {
		t.Error("setting with login token, should be passed")
	}

	settingDNSPod = &Settings{Provider: "DNSPod"}
	if err := CheckSettings(settingDNSPod); err == nil {
		t.Error("setting with invalid parameters, should be failed")
	}

	settingHE := &Settings{Provider: "HE", Password: ""}
	if err := CheckSettings(settingHE); err != nil {
		t.Log("HE setting without password, passed")
	} else {
		t.Error("HE setting without password, should be faild")
	}
}
