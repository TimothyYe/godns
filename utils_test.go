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

func TestSaveCurrentIP(t *testing.T) {
	SaveCurrentIP("1.2.3.4")

	if _, err := os.Stat("./.current_ip"); os.IsNotExist(err) {
		t.Error(".current_ip file should exists")
	}

	savedIP := LoadCurrentIP()

	if strings.TrimRight(savedIP, "\n") != "1.2.3.4" {
		t.Error("saved IP should be equal to 1.2.3.4")
	}

	//Cleanup
	os.Remove("./.current_ip")
}

func TestLoadCurrentIP(t *testing.T) {
	ip := LoadCurrentIP()

	if ip != "" {
		t.Error("current ip file should be empth")
	}
}
