package provider

import (
	"testing"

	"github.com/TimothyYe/godns/internal/settings"
	"github.com/TimothyYe/godns/internal/utils"
)

func TestIntegrationSingleProvider(t *testing.T) {
	t.Run("DNSPodSingleProvider", func(t *testing.T) {
		config := &settings.Settings{
			Provider:   "DNSPod",
			LoginToken: "test-token",
			Domains: []settings.Domain{
				{DomainName: "example.com", SubDomains: []string{"www", "api"}},
				{DomainName: "example.net", SubDomains: []string{"mail", "ftp"}},
			},
		}

		// Validate configuration
		if err := utils.CheckSettings(config); err != nil {
			t.Fatalf("Configuration validation failed: %v", err)
		}

		// Test provider creation
		providers, err := GetProviders(config)
		if err != nil {
			t.Fatalf("Failed to create providers: %v", err)
		}

		if len(providers) != 1 {
			t.Errorf("Expected 1 provider, got %d", len(providers))
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
	})

	t.Run("CloudflareSingleProvider", func(t *testing.T) {
		config := &settings.Settings{
			Provider: "Cloudflare",
			Email:    "test@example.com",
			Password: "test-token",
			Domains: []settings.Domain{
				{DomainName: "cloudflare-domain.com", SubDomains: []string{"www", "api", "@"}},
			},
		}

		// Validate configuration
		if err := utils.CheckSettings(config); err != nil {
			t.Fatalf("Configuration validation failed: %v", err)
		}

		// Test provider creation
		providers, err := GetProviders(config)
		if err != nil {
			t.Fatalf("Failed to create providers: %v", err)
		}

		if len(providers) != 1 {
			t.Errorf("Expected 1 provider, got %d", len(providers))
		}

		// Test domain resolution
		domain := &config.Domains[0]
		providerName := config.GetDomainProvider(domain)
		if providerName != "Cloudflare" {
			t.Errorf("Expected Cloudflare for domain %s, got %s", domain.DomainName, providerName)
		}

		provider, err := GetProviderForDomain(domain, providers, config)
		if err != nil {
			t.Errorf("Failed to get provider for domain %s: %v", domain.DomainName, err)
		}
		if provider == nil {
			t.Errorf("Provider is nil for domain %s", domain.DomainName)
		}
	})
}

func TestIntegrationMultiProvider(t *testing.T) {
	t.Run("FullMultiProvider", func(t *testing.T) {
		config := &settings.Settings{
			Providers: map[string]*settings.ProviderConfig{
				"Cloudflare": {
					Email:    "test@example.com",
					Password: "cf-token",
				},
				"DNSPod": {
					LoginToken: "dnspod-token",
				},
				"DigitalOcean": {
					LoginToken: "do-token",
				},
				"DuckDNS": {
					LoginToken: "duck-token",
				},
			},
			Domains: []settings.Domain{
				{DomainName: "cf-domain.com", SubDomains: []string{"www", "api"}, Provider: "Cloudflare"},
				{DomainName: "dnspod-domain.com", SubDomains: []string{"www", "mail"}, Provider: "DNSPod"},
				{DomainName: "do-domain.com", SubDomains: []string{"www", "blog"}, Provider: "DigitalOcean"},
				{DomainName: "duck-domain.org", SubDomains: []string{"home"}, Provider: "DuckDNS"},
			},
		}

		// Validate configuration
		if err := utils.CheckSettings(config); err != nil {
			t.Fatalf("Configuration validation failed: %v", err)
		}

		// Test provider creation
		providers, err := GetProviders(config)
		if err != nil {
			t.Fatalf("Failed to create providers: %v", err)
		}

		if len(providers) != 4 {
			t.Errorf("Expected 4 providers, got %d", len(providers))
		}

		// Test each domain resolution
		expectedMappings := map[string]string{
			"cf-domain.com":     "Cloudflare",
			"dnspod-domain.com": "DNSPod",
			"do-domain.com":     "DigitalOcean",
			"duck-domain.org":   "DuckDNS",
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
	})

	t.Run("MultiProviderWithSpecialProviders", func(t *testing.T) {
		config := &settings.Settings{
			Providers: map[string]*settings.ProviderConfig{
				"OVH": {
					AppKey:      "ovh-app-key",
					AppSecret:   "ovh-app-secret",
					ConsumerKey: "ovh-consumer-key",
				},
				"AliDNS": {
					Email:    "test@ali.com",
					Password: "ali-password",
				},
				"HE": {
					Password: "he-password",
				},
			},
			Domains: []settings.Domain{
				{DomainName: "ovh-domain.fr", SubDomains: []string{"www"}, Provider: "OVH"},
				{DomainName: "ali-domain.cn", SubDomains: []string{"www"}, Provider: "AliDNS"},
				{DomainName: "he-domain.net", SubDomains: []string{"www"}, Provider: "HE"},
			},
		}

		// Validate configuration
		if err := utils.CheckSettings(config); err != nil {
			t.Fatalf("Configuration validation failed: %v", err)
		}

		// Test provider creation
		providers, err := GetProviders(config)
		if err != nil {
			t.Fatalf("Failed to create providers: %v", err)
		}

		if len(providers) != 3 {
			t.Errorf("Expected 3 providers, got %d", len(providers))
		}

		// Verify all providers exist
		expectedProviders := []string{"OVH", "AliDNS", "HE"}
		for _, expectedProvider := range expectedProviders {
			if _, exists := providers[expectedProvider]; !exists {
				t.Errorf("Expected provider %s not found", expectedProvider)
			}
		}
	})
}

func TestIntegrationMixedProvider(t *testing.T) {
	t.Run("MixedDynuDuckDNS", func(t *testing.T) {
		// This mirrors your actual config.json structure
		config := &settings.Settings{
			Providers: map[string]*settings.ProviderConfig{
				"DuckDNS": {
					LoginToken: "duck-token",
				},
				"Dynu": {
					Password: "dynu-password",
				},
			},
			Domains: []settings.Domain{
				{DomainName: "ddnsfree.com", SubDomains: []string{"godns"}, Provider: "Dynu"},
				{DomainName: "duckdns.org", SubDomains: []string{"itimothyye"}, Provider: "DuckDNS"},
			},
		}

		// Validate configuration
		if err := utils.CheckSettings(config); err != nil {
			t.Fatalf("Configuration validation failed: %v", err)
		}

		// Test provider creation
		providers, err := GetProviders(config)
		if err != nil {
			t.Fatalf("Failed to create providers: %v", err)
		}

		if len(providers) != 2 {
			t.Errorf("Expected 2 providers, got %d", len(providers))
		}

		// Verify both providers exist
		if _, exists := providers["Dynu"]; !exists {
			t.Error("Dynu provider not found")
		}
		if _, exists := providers["DuckDNS"]; !exists {
			t.Error("DuckDNS provider not found")
		}

		// Test domain resolution
		for _, domain := range config.Domains {
			providerName := config.GetDomainProvider(&domain)
			expectedProvider := domain.Provider

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
	})

	t.Run("MixedGlobalPlusSpecific", func(t *testing.T) {
		config := &settings.Settings{
			Provider:   "DNSPod",
			LoginToken: "global-dnspod-token",
			Providers: map[string]*settings.ProviderConfig{
				"Cloudflare": {
					LoginToken: "cf-token",
				},
				"DigitalOcean": {
					LoginToken: "do-token",
				},
			},
			Domains: []settings.Domain{
				{DomainName: "old-domain.com", SubDomains: []string{"www", "mail"}}, // Uses global DNSPod
				{DomainName: "cf-domain.com", SubDomains: []string{"www"}, Provider: "Cloudflare"},
				{DomainName: "do-domain.com", SubDomains: []string{"www"}, Provider: "DigitalOcean"},
			},
		}

		// Validate configuration
		if err := utils.CheckSettings(config); err != nil {
			t.Fatalf("Configuration validation failed: %v", err)
		}

		// Test provider creation
		providers, err := GetProviders(config)
		if err != nil {
			t.Fatalf("Failed to create providers: %v", err)
		}

		if len(providers) != 3 {
			t.Errorf("Expected 3 providers, got %d", len(providers))
		}

		// Verify all providers exist including global
		expectedProviders := []string{"DNSPod", "Cloudflare", "DigitalOcean"}
		for _, expectedProvider := range expectedProviders {
			if _, exists := providers[expectedProvider]; !exists {
				t.Errorf("Expected provider %s not found", expectedProvider)
			}
		}

		// Test domain resolution
		expectedMappings := map[string]string{
			"old-domain.com": "DNSPod", // Uses global provider
			"cf-domain.com":  "Cloudflare",
			"do-domain.com":  "DigitalOcean",
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
	})

	t.Run("MixedGlobalTakesPrecedence", func(t *testing.T) {
		// Test that global provider takes precedence when same name exists in providers
		config := &settings.Settings{
			Provider: "Cloudflare",
			Email:    "global@example.com",
			Password: "global-password",
			Providers: map[string]*settings.ProviderConfig{
				"Cloudflare": {
					Email:    "specific@example.com",
					Password: "specific-password",
				},
				"DNSPod": {
					LoginToken: "dnspod-token",
				},
			},
			Domains: []settings.Domain{
				{DomainName: "global-cf.com", SubDomains: []string{"www"}}, // Should use global Cloudflare
				{DomainName: "dnspod.com", SubDomains: []string{"www"}, Provider: "DNSPod"},
			},
		}

		// Validate configuration
		if err := utils.CheckSettings(config); err != nil {
			t.Fatalf("Configuration validation failed: %v", err)
		}

		// Test provider creation
		providers, err := GetProviders(config)
		if err != nil {
			t.Fatalf("Failed to create providers: %v", err)
		}

		// Should have 2 providers (global Cloudflare takes precedence, plus DNSPod)
		if len(providers) != 2 {
			t.Errorf("Expected 2 providers, got %d", len(providers))
		}

		// Verify providers exist
		if _, exists := providers["Cloudflare"]; !exists {
			t.Error("Cloudflare provider not found")
		}
		if _, exists := providers["DNSPod"]; !exists {
			t.Error("DNSPod provider not found")
		}

		// Test domain resolution
		for _, domain := range config.Domains {
			providerName := config.GetDomainProvider(&domain)

			provider, err := GetProviderForDomain(&domain, providers, config)
			if err != nil {
				t.Errorf("Failed to get provider for domain %s: %v", domain.DomainName, err)
			}
			if provider == nil {
				t.Errorf("Provider is nil for domain %s", domain.DomainName)
			}

			// Verify correct provider assignment
			if domain.DomainName == "global-cf.com" && providerName != "Cloudflare" {
				t.Errorf("Expected Cloudflare for global-cf.com, got %s", providerName)
			}
			if domain.DomainName == "dnspod.com" && providerName != "DNSPod" {
				t.Errorf("Expected DNSPod for dnspod.com, got %s", providerName)
			}
		}
	})
}

func TestIntegrationErrorCases(t *testing.T) {
	t.Run("InvalidProviderReference", func(t *testing.T) {
		config := &settings.Settings{
			Providers: map[string]*settings.ProviderConfig{
				"Cloudflare": {
					LoginToken: "cf-token",
				},
			},
			Domains: []settings.Domain{
				{DomainName: "example.com", SubDomains: []string{"www"}, Provider: "NonExistentProvider"},
			},
		}

		// Configuration validation should fail
		if err := utils.CheckSettings(config); err == nil {
			t.Error("Expected configuration validation to fail, but it passed")
		}
	})

	t.Run("MissingProviderCredentials", func(t *testing.T) {
		config := &settings.Settings{
			Providers: map[string]*settings.ProviderConfig{
				"Cloudflare": {
					// Missing credentials
				},
			},
			Domains: []settings.Domain{
				{DomainName: "example.com", SubDomains: []string{"www"}, Provider: "Cloudflare"},
			},
		}

		// Configuration validation should fail
		if err := utils.CheckSettings(config); err == nil {
			t.Error("Expected configuration validation to fail for missing credentials, but it passed")
		}
	})

	t.Run("EmptyDomainName", func(t *testing.T) {
		config := &settings.Settings{
			Provider:   "DNSPod",
			LoginToken: "test-token",
			Domains: []settings.Domain{
				{DomainName: "", SubDomains: []string{"www"}}, // Empty domain name
			},
		}

		// Configuration validation should fail
		if err := utils.CheckSettings(config); err == nil {
			t.Error("Expected configuration validation to fail for empty domain name, but it passed")
		}
	})

	t.Run("EmptySubdomain", func(t *testing.T) {
		config := &settings.Settings{
			Provider:   "DNSPod",
			LoginToken: "test-token",
			Domains: []settings.Domain{
				{DomainName: "example.com", SubDomains: []string{""}}, // Empty subdomain
			},
		}

		// Configuration validation should fail
		if err := utils.CheckSettings(config); err == nil {
			t.Error("Expected configuration validation to fail for empty subdomain, but it passed")
		}
	})
}

// Benchmark tests for performance.
func BenchmarkGetProviders(b *testing.B) {
	config := &settings.Settings{
		Providers: map[string]*settings.ProviderConfig{
			"Cloudflare": {
				LoginToken: "cf-token",
			},
			"DNSPod": {
				LoginToken: "dnspod-token",
			},
			"DigitalOcean": {
				LoginToken: "do-token",
			},
		},
		Domains: []settings.Domain{
			{DomainName: "example.com", SubDomains: []string{"www"}, Provider: "Cloudflare"},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := GetProviders(config)
		if err != nil {
			b.Fatalf("GetProviders failed: %v", err)
		}
	}
}

func BenchmarkGetProviderForDomain(b *testing.B) {
	config := &settings.Settings{
		Providers: map[string]*settings.ProviderConfig{
			"Cloudflare": {
				LoginToken: "cf-token",
			},
		},
	}

	providers, err := GetProviders(config)
	if err != nil {
		b.Fatalf("Failed to create providers: %v", err)
	}

	domain := &settings.Domain{
		DomainName: "example.com",
		SubDomains: []string{"www"},
		Provider:   "Cloudflare",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := GetProviderForDomain(domain, providers, config)
		if err != nil {
			b.Fatalf("GetProviderForDomain failed: %v", err)
		}
	}
}
