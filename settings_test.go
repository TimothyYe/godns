package main

import (
	"testing"
)

func TestLoadSetting(t *testing.T) {
	settings, err := LoadSettings("./config_sample.json")

	if err != nil {
		t.Error(err.Error())
	}

	if settings.IP_Url == "" {
		t.Error("Cannot load ip_url from config file")
	}
}
