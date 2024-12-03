package utils

import (
	"testing"

	"github.com/TimothyYe/godns/internal/settings"
)

func TestCheckSettings(t *testing.T) {
	settingError := &settings.Settings{}
	if err := CheckSettings(settingError); err == nil {
		t.Error("setting is invalid, should return error")
	}

	settingDNSPod := &settings.Settings{Provider: "DNSPod", LoginToken: "aaa"}
	if err := CheckSettings(settingDNSPod); err == nil {
		t.Log("setting with login token, passed")
	} else {
		t.Error("setting with login token, should be passed")
	}

	settingDNSPod = &settings.Settings{Provider: "DNSPod"}
	if err := CheckSettings(settingDNSPod); err == nil {
		t.Error("setting with invalid parameters, should be failed")
	}

	settingHE := &settings.Settings{Provider: "HE", Password: ""}
	if err := CheckSettings(settingHE); err != nil {
		t.Log("HE setting without password, passed")
	} else {
		t.Error("HE setting without password, should be failed")
	}
}
