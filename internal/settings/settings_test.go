package settings

import (
	"testing"
)

func TestLoadSetting(t *testing.T) {
	var settings Settings
	err := LoadSettings("../../configs/config_sample.json", &settings)

	if err != nil {
		t.Fatal(err.Error())
	}

	if settings.IPUrl == "" {
		t.Fatal("cannot load ip_url from config file")
	}

	err = LoadSettings("./file/does/not/exists", &settings)
	if err == nil {
		t.Fatal("file doesn't exist, should return error")
	}
}
