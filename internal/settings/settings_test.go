package settings

import (
	"testing"
)

func TestLoadJSONSetting(t *testing.T) {
	var settings Settings
	err := LoadSettings("../../configs/config_sample.json", &settings)

	if err != nil {
		t.Fatal(err.Error())
	}

	if len(settings.IPUrls) == 0 && settings.IPUrl == "" {
		t.Fatal("neither ip_urls nor ip_url contain valid entries")
	}

	err = LoadSettings("./file/does/not/exists", &settings)
	if err == nil {
		t.Fatal("file doesn't exist, should return error")
	}
}

func TestLoadYAMLSetting(t *testing.T) {
	var settings Settings
	err := LoadSettings("../../configs/config_sample.yaml", &settings)

	if err != nil {
		t.Fatal(err.Error())
	}

	if len(settings.IPUrls) == 0 && settings.IPUrl == "" {
		t.Fatal("cannot load ip_url from config file")
	}

	t.Log(settings)
}
