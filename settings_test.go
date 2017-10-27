package godns

import (
	"testing"
)

func TestLoadSetting(t *testing.T) {
	var settings Settings
	err := LoadSettings("./config_sample.json", &settings)

	if err != nil {
		t.Error(err.Error())
	}

	if settings.IPUrl == "" {
		t.Error("Cannot load ip_url from config file")
	}
}
