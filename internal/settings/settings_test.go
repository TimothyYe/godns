package settings

import (
	"os"
	"path/filepath"
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

func TestLoadWithEnvPath(t *testing.T) {
	const expectedToken = "super secret login token"

	tokenFile, err := os.CreateTemp("", "login_token_*.txt")
	if err != nil {
		t.Fatalf("failed to create temp token file: %v", err)
	}
	defer os.Remove(tokenFile.Name())

	if _, err := tokenFile.WriteString(expectedToken); err != nil {
		t.Fatalf("failed to write token to file: %v", err)
	}
	if err := tokenFile.Close(); err != nil {
		t.Fatalf("failed to close token file: %v", err)
	}

	settings := Settings{
		LoginTokenFile: "$LOGIN_TOKEN_FILE",
	}

	configFile, err := os.CreateTemp("", "config_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp config file: %v", err)
	}
	defer os.Remove(configFile.Name())

	if err := settings.SaveSettings(configFile.Name()); err != nil {
		t.Fatalf("failed to save settings to config file: %v", err)
	}

	t.Setenv("LOGIN_TOKEN_FILE", tokenFile.Name())

	var loaded Settings
	if err := LoadSettings(configFile.Name(), &loaded); err != nil {
		t.Fatalf("failed to load settings with env var: %v", err)
	}

	if loaded.LoginToken != expectedToken {
		t.Errorf("expected login token %q, got %q", expectedToken, loaded.LoginToken)
	}

	t.Log(settings)
}

func TestQueryInterfaceField(t *testing.T) {
	dir := t.TempDir()

	jsonPath := filepath.Join(dir, "config.json")
	if err := os.WriteFile(jsonPath, []byte(`{"query_interface": "wan0"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	var jsonSettings Settings
	if err := LoadSettings(jsonPath, &jsonSettings); err != nil {
		t.Fatal(err)
	}
	if jsonSettings.QueryInterface != "wan0" {
		t.Errorf("json query_interface = %q, want wan0", jsonSettings.QueryInterface)
	}

	yamlPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(yamlPath, []byte("query_interface: wan1\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	var yamlSettings Settings
	if err := LoadSettings(yamlPath, &yamlSettings); err != nil {
		t.Fatal(err)
	}
	if yamlSettings.QueryInterface != "wan1" {
		t.Errorf("yaml query_interface = %q, want wan1", yamlSettings.QueryInterface)
	}
}
