package provider

import (
	"path/filepath"
	"testing"

	"github.com/TimothyYe/godns/internal/settings"
	"github.com/TimothyYe/godns/internal/utils"
)

func TestConfigurationFiles(t *testing.T) {
	testDataDir := filepath.Join("..", "..", "testdata")

	t.Run("SingleProviderConfigFile", func(t *testing.T) {
		configPath := filepath.Join(testDataDir, "config_test_single.json")
		config := &settings.Settings{}

		// Load configuration from file
		err := settings.LoadSettings(configPath, config)
		if err != nil {
			t.Fatalf("Failed to load single provider config: %v", err)
		}

		// Validate configuration
		if err := utils.CheckSettings(config); err != nil {
			t.Fatalf("Single provider configuration validation failed: %v", err)
		}

		// Test provider creation
		providers, err := GetProviders(config)
		if err != nil {
			t.Fatalf("Failed to create providers from single config: %v", err)
		}

		if len(providers) != 1 {
			t.Errorf("Expected 1 provider for single config, got %d", len(providers))
		}

		if _, exists := providers["DNSPod"]; !exists {
			t.Error("DNSPod provider not found in single config")
		}

		// Test domain resolution
		for _, domain := range config.Domains {
			providerName := config.GetDomainProvider(&domain)
			if providerName != "DNSPod" {
				t.Errorf("Expected DNSPod for domain %s, got %s", domain.DomainName, providerName)
			}

			provider, err := GetProviderForDomain(&domain, providers, config)
			if err != nil {
				t.Errorf("Failed to get provider for domain %s: %v", domain.DomainName, err)
			}
			if provider == nil {
				t.Errorf("Provider is nil for domain %s", domain.DomainName)
			}
		}

		t.Logf("Single provider config test passed with %d domains", len(config.Domains))
	})

	t.Run("MultiProviderConfigFile", func(t *testing.T) {
		configPath := filepath.Join(testDataDir, "config_test_multi.json")
		config := &settings.Settings{}

		// Load configuration from file
		err := settings.LoadSettings(configPath, config)
		if err != nil {
			t.Fatalf("Failed to load multi provider config: %v", err)
		}

		// Validate configuration
		if err := utils.CheckSettings(config); err != nil {
			t.Fatalf("Multi provider configuration validation failed: %v", err)
		}

		// Test provider creation
		providers, err := GetProviders(config)
		if err != nil {
			t.Fatalf("Failed to create providers from multi config: %v", err)
		}

		if len(providers) != 4 {
			t.Errorf("Expected 4 providers for multi config, got %d", len(providers))
		}

		expectedProviders := []string{"Cloudflare", "DNSPod", "DigitalOcean", "DuckDNS"}
		for _, expectedProvider := range expectedProviders {
			if _, exists := providers[expectedProvider]; !exists {
				t.Errorf("%s provider not found in multi config", expectedProvider)
			}
		}

		// Test domain resolution with expected mappings
		expectedMappings := map[string]string{
			"cloudflare-test.com":   "Cloudflare",
			"dnspod-test.com":       "DNSPod",
			"digitalocean-test.com": "DigitalOcean",
			"duckdns-test.org":      "DuckDNS",
		}

		for _, domain := range config.Domains {
			expectedProvider := expectedMappings[domain.DomainName]
			providerName := config.GetDomainProvider(&domain)

			if providerName != expectedProvider {
				t.Errorf("Expected %s for domain %s, got %s", expectedProvider, domain.DomainName, providerName)
			}

			provider, err := GetProviderForDomain(&domain, providers, config)
			if err != nil {
				t.Errorf("Failed to get provider for domain %s: %v", domain.DomainName, err)
			}
			if provider == nil {
				t.Errorf("Provider is nil for domain %s", domain.DomainName)
			}
		}

		t.Logf("Multi provider config test passed with %d providers and %d domains", len(providers), len(config.Domains))
	})

	t.Run("MixedProviderConfigFile", func(t *testing.T) {
		configPath := filepath.Join(testDataDir, "config_test_mixed.json")
		config := &settings.Settings{}

		// Load configuration from file
		err := settings.LoadSettings(configPath, config)
		if err != nil {
			t.Fatalf("Failed to load mixed provider config: %v", err)
		}

		// Validate configuration
		if err := utils.CheckSettings(config); err != nil {
			t.Fatalf("Mixed provider configuration validation failed: %v", err)
		}

		// Test provider creation
		providers, err := GetProviders(config)
		if err != nil {
			t.Fatalf("Failed to create providers from mixed config: %v", err)
		}

		if len(providers) != 3 {
			t.Errorf("Expected 3 providers for mixed config, got %d", len(providers))
		}

		// Should have global DNSPod + specific providers
		expectedProviders := []string{"DNSPod", "Cloudflare", "DigitalOcean"}
		for _, expectedProvider := range expectedProviders {
			if _, exists := providers[expectedProvider]; !exists {
				t.Errorf("%s provider not found in mixed config", expectedProvider)
			}
		}

		// Test domain resolution with expected mappings
		expectedMappings := map[string]string{
			"legacy-domain.com":      "DNSPod",       // Uses global provider
			"cloudflare-mixed.com":   "Cloudflare",   // Uses specific provider
			"digitalocean-mixed.com": "DigitalOcean", // Uses specific provider
			"another-legacy.net":     "DNSPod",       // Uses global provider
		}

		for _, domain := range config.Domains {
			expectedProvider := expectedMappings[domain.DomainName]
			providerName := config.GetDomainProvider(&domain)

			if providerName != expectedProvider {
				t.Errorf("Expected %s for domain %s, got %s", expectedProvider, domain.DomainName, providerName)
			}

			provider, err := GetProviderForDomain(&domain, providers, config)
			if err != nil {
				t.Errorf("Failed to get provider for domain %s: %v", domain.DomainName, err)
			}
			if provider == nil {
				t.Errorf("Provider is nil for domain %s", domain.DomainName)
			}
		}

		t.Logf("Mixed provider config test passed with %d providers and %d domains", len(providers), len(config.Domains))
	})
}

func TestConfigurationModeDetection(t *testing.T) {
	testDataDir := filepath.Join("..", "..", "testdata")

	testCases := []struct {
		name          string
		configFile    string
		expectedMode  string
		expectedMulti bool
	}{
		{
			name:          "Single Provider Mode",
			configFile:    "config_test_single.json",
			expectedMode:  "single",
			expectedMulti: false,
		},
		{
			name:          "Multi Provider Mode",
			configFile:    "config_test_multi.json",
			expectedMode:  "multi",
			expectedMulti: true,
		},
		{
			name:          "Mixed Provider Mode",
			configFile:    "config_test_mixed.json",
			expectedMode:  "mixed",
			expectedMulti: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			configPath := filepath.Join(testDataDir, tc.configFile)
			config := &settings.Settings{}

			err := settings.LoadSettings(configPath, config)
			if err != nil {
				t.Fatalf("Failed to load config %s: %v", tc.configFile, err)
			}

			isMulti := config.IsMultiProvider()
			if isMulti != tc.expectedMulti {
				t.Errorf("Expected IsMultiProvider() to be %v for %s, got %v", tc.expectedMulti, tc.name, isMulti)
			}

			// Additional checks based on mode
			switch tc.expectedMode {
			case "single":
				if config.Provider == "" {
					t.Error("Single provider mode should have Provider field set")
				}
				if len(config.Providers) > 0 {
					t.Error("Single provider mode should not have Providers map")
				}

			case "multi":
				if config.Provider != "" {
					t.Error("Pure multi provider mode should not have global Provider field set")
				}
				if len(config.Providers) == 0 {
					t.Error("Multi provider mode should have Providers map")
				}

			case "mixed":
				if config.Provider == "" {
					t.Error("Mixed provider mode should have global Provider field set")
				}
				if len(config.Providers) == 0 {
					t.Error("Mixed provider mode should have Providers map")
				}
			}

			t.Logf("%s detected correctly: IsMultiProvider=%v", tc.name, isMulti)
		})
	}
}

func TestRealWorldScenarios(t *testing.T) {
	t.Run("MigrationScenario", func(t *testing.T) {
		// Simulate a user migrating from single to mixed to multi provider

		// Step 1: Start with single provider
		singleConfig := &settings.Settings{
			Provider:   "DNSPod",
			LoginToken: "dnspod-token",
			Domains: []settings.Domain{
				{DomainName: "example.com", SubDomains: []string{"www"}},
			},
		}

		if err := utils.CheckSettings(singleConfig); err != nil {
			t.Fatalf("Single config validation failed: %v", err)
		}

		providers1, err := GetProviders(singleConfig)
		if err != nil {
			t.Fatalf("Single config provider creation failed: %v", err)
		}
		if len(providers1) != 1 {
			t.Errorf("Expected 1 provider in single mode, got %d", len(providers1))
		}

		// Step 2: Add new provider while keeping old (mixed mode)
		mixedConfig := &settings.Settings{
			Provider:   "DNSPod",
			LoginToken: "dnspod-token",
			Providers: map[string]*settings.ProviderConfig{
				"Cloudflare": {
					LoginToken: "cf-token",
				},
			},
			Domains: []settings.Domain{
				{DomainName: "example.com", SubDomains: []string{"www"}}, // Uses global DNSPod
				{DomainName: "newsite.com", SubDomains: []string{"www"}, Provider: "Cloudflare"},
			},
		}

		if err := utils.CheckSettings(mixedConfig); err != nil {
			t.Fatalf("Mixed config validation failed: %v", err)
		}

		providers2, err := GetProviders(mixedConfig)
		if err != nil {
			t.Fatalf("Mixed config provider creation failed: %v", err)
		}
		if len(providers2) != 2 {
			t.Errorf("Expected 2 providers in mixed mode, got %d", len(providers2))
		}

		// Step 3: Convert to full multi-provider mode
		multiConfig := &settings.Settings{
			Providers: map[string]*settings.ProviderConfig{
				"DNSPod": {
					LoginToken: "dnspod-token",
				},
				"Cloudflare": {
					LoginToken: "cf-token",
				},
			},
			Domains: []settings.Domain{
				{DomainName: "example.com", SubDomains: []string{"www"}, Provider: "DNSPod"},
				{DomainName: "newsite.com", SubDomains: []string{"www"}, Provider: "Cloudflare"},
			},
		}

		if err := utils.CheckSettings(multiConfig); err != nil {
			t.Fatalf("Multi config validation failed: %v", err)
		}

		providers3, err := GetProviders(multiConfig)
		if err != nil {
			t.Fatalf("Multi config provider creation failed: %v", err)
		}
		if len(providers3) != 2 {
			t.Errorf("Expected 2 providers in multi mode, got %d", len(providers3))
		}

		t.Log("Migration scenario test passed: single -> mixed -> multi")
	})

	t.Run("LargeScaleDeployment", func(t *testing.T) {
		// Test a large-scale deployment with many providers and domains
		config := &settings.Settings{
			Providers: map[string]*settings.ProviderConfig{
				"Cloudflare":   {LoginToken: "cf-token"},
				"DNSPod":       {LoginToken: "dnspod-token"},
				"DigitalOcean": {LoginToken: "do-token"},
				"DuckDNS":      {LoginToken: "duck-token"},
				"Google":       {Email: "test@gmail.com", Password: "google-pass"},
				"AliDNS":       {Email: "test@ali.com", Password: "ali-pass"},
			},
			Domains: []settings.Domain{
				{DomainName: "site1.com", SubDomains: []string{"www", "api", "mail"}, Provider: "Cloudflare"},
				{DomainName: "site2.com", SubDomains: []string{"www", "blog"}, Provider: "DNSPod"},
				{DomainName: "site3.org", SubDomains: []string{"www"}, Provider: "DigitalOcean"},
				{DomainName: "home.duckdns.org", SubDomains: []string{"myhouse"}, Provider: "DuckDNS"},
				{DomainName: "corp.com", SubDomains: []string{"www", "intranet"}, Provider: "Google"},
				{DomainName: "china.cn", SubDomains: []string{"www", "api"}, Provider: "AliDNS"},
			},
		}

		if err := utils.CheckSettings(config); err != nil {
			t.Fatalf("Large scale config validation failed: %v", err)
		}

		providers, err := GetProviders(config)
		if err != nil {
			t.Fatalf("Large scale provider creation failed: %v", err)
		}

		if len(providers) != 6 {
			t.Errorf("Expected 6 providers in large scale deployment, got %d", len(providers))
		}

		// Test all domains resolve correctly
		for _, domain := range config.Domains {
			providerName := config.GetDomainProvider(&domain)
			if providerName != domain.Provider {
				t.Errorf("Domain %s: expected provider %s, got %s", domain.DomainName, domain.Provider, providerName)
			}

			provider, err := GetProviderForDomain(&domain, providers, config)
			if err != nil {
				t.Errorf("Failed to get provider for domain %s: %v", domain.DomainName, err)
			}
			if provider == nil {
				t.Errorf("Provider is nil for domain %s", domain.DomainName)
			}
		}

		t.Logf("Large scale deployment test passed with %d providers and %d domains", len(providers), len(config.Domains))
	})
}
