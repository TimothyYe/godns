package provider

import (
	"testing"

	"github.com/TimothyYe/godns/internal/settings"
)

func TestGetProviders(t *testing.T) {
	// Test single provider mode
	t.Run("SingleProviderMode", func(t *testing.T) {
		// Test DNSPod single provider
		t.Run("DNSPod", func(t *testing.T) {
			config := &settings.Settings{
				Provider:   "DNSPod",
				LoginToken: "test-token",
				Domains: []settings.Domain{
					{DomainName: "example.com", SubDomains: []string{"www"}},
				},
			}

			providers, err := GetProviders(config)
			if err != nil {
				t.Fatalf("Failed to get providers: %v", err)
			}

			if len(providers) != 1 {
				t.Errorf("Expected 1 provider, got %d", len(providers))
			}

			if _, exists := providers["DNSPod"]; !exists {
				t.Error("DNSPod provider not found")
			}
		})

		// Test Cloudflare single provider
		t.Run("Cloudflare", func(t *testing.T) {
			config := &settings.Settings{
				Provider:   "Cloudflare",
				LoginToken: "test-token",
				Domains: []settings.Domain{
					{DomainName: "example.com", SubDomains: []string{"www"}},
				},
			}

			providers, err := GetProviders(config)
			if err != nil {
				t.Fatalf("Failed to get providers: %v", err)
			}

			if len(providers) != 1 {
				t.Errorf("Expected 1 provider, got %d", len(providers))
			}

			if _, exists := providers["Cloudflare"]; !exists {
				t.Error("Cloudflare provider not found")
			}
		})

		// Test empty provider should fail
		t.Run("EmptyProvider", func(t *testing.T) {
			config := &settings.Settings{
				Provider: "",
				Domains: []settings.Domain{
					{DomainName: "example.com", SubDomains: []string{"www"}},
				},
			}

			_, err := GetProviders(config)
			if err == nil {
				t.Error("Expected error for empty provider, but got none")
			}
		})
	})

	// Test multi-provider mode
	t.Run("MultiProviderMode", func(t *testing.T) {
		// Test multiple providers
		t.Run("MultipleProviders", func(t *testing.T) {
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
				},
				Domains: []settings.Domain{
					{DomainName: "example.com", SubDomains: []string{"www"}, Provider: "Cloudflare"},
					{DomainName: "example.net", SubDomains: []string{"www"}, Provider: "DNSPod"},
					{DomainName: "example.org", SubDomains: []string{"www"}, Provider: "DigitalOcean"},
				},
			}

			providers, err := GetProviders(config)
			if err != nil {
				t.Fatalf("Failed to get providers: %v", err)
			}

			if len(providers) != 3 {
				t.Errorf("Expected 3 providers, got %d", len(providers))
			}

			expectedProviders := []string{"Cloudflare", "DNSPod", "DigitalOcean"}
			for _, providerName := range expectedProviders {
				if _, exists := providers[providerName]; !exists {
					t.Errorf("%s provider not found", providerName)
				}
			}
		})

		// Test OVH with special keys
		t.Run("OVHProvider", func(t *testing.T) {
			config := &settings.Settings{
				Providers: map[string]*settings.ProviderConfig{
					"OVH": {
						AppKey:      "test-app-key",
						AppSecret:   "test-app-secret",
						ConsumerKey: "test-consumer-key",
					},
				},
				Domains: []settings.Domain{
					{DomainName: "example.com", SubDomains: []string{"www"}, Provider: "OVH"},
				},
			}

			providers, err := GetProviders(config)
			if err != nil {
				t.Fatalf("Failed to get providers: %v", err)
			}

			if len(providers) != 1 {
				t.Errorf("Expected 1 provider, got %d", len(providers))
			}

			if _, exists := providers["OVH"]; !exists {
				t.Error("OVH provider not found")
			}
		})

		// Test empty providers map with no global provider - should fail
		t.Run("EmptyProviders", func(t *testing.T) {
			config := &settings.Settings{
				Providers: map[string]*settings.ProviderConfig{},
				Domains: []settings.Domain{
					{DomainName: "example.com", SubDomains: []string{"www"}},
				},
			}

			_, err := GetProviders(config)
			if err == nil {
				t.Error("Expected error for configuration with no providers, but got none")
			}
		})
	})

	// Test mixed configuration mode
	t.Run("MixedConfigurationMode", func(t *testing.T) {
		// Test global provider + specific providers
		t.Run("GlobalPlusSpecific", func(t *testing.T) {
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
					{DomainName: "old-domain.com", SubDomains: []string{"www"}}, // Uses global DNSPod
					{DomainName: "new-domain.com", SubDomains: []string{"www"}, Provider: "Cloudflare"},
					{DomainName: "another-domain.com", SubDomains: []string{"www"}, Provider: "DigitalOcean"},
				},
			}

			providers, err := GetProviders(config)
			if err != nil {
				t.Fatalf("Failed to get providers: %v", err)
			}

			if len(providers) != 3 {
				t.Errorf("Expected 3 providers, got %d", len(providers))
			}

			expectedProviders := []string{"DNSPod", "Cloudflare", "DigitalOcean"}
			for _, providerName := range expectedProviders {
				if _, exists := providers[providerName]; !exists {
					t.Errorf("%s provider not found", providerName)
				}
			}
		})

		// Test global provider takes precedence over same name in providers
		t.Run("GlobalTakesPrecedence", func(t *testing.T) {
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
					{DomainName: "example.com", SubDomains: []string{"www"}}, // Should use global Cloudflare
					{DomainName: "example.net", SubDomains: []string{"www"}, Provider: "DNSPod"},
				},
			}

			providers, err := GetProviders(config)
			if err != nil {
				t.Fatalf("Failed to get providers: %v", err)
			}

			if len(providers) != 2 {
				t.Errorf("Expected 2 providers, got %d", len(providers))
			}

			if _, exists := providers["Cloudflare"]; !exists {
				t.Error("Cloudflare provider not found")
			}
			if _, exists := providers["DNSPod"]; !exists {
				t.Error("DNSPod provider not found")
			}
		})

		// Test mixed with various provider types
		t.Run("MixedVariousProviders", func(t *testing.T) {
			config := &settings.Settings{
				Provider: "HE",
				Password: "he-password",
				Providers: map[string]*settings.ProviderConfig{
					"DuckDNS": {
						LoginToken: "duck-token",
					},
					"Google": {
						Email:    "test@gmail.com",
						Password: "google-password",
					},
					"AliDNS": {
						Email:    "test@ali.com",
						Password: "ali-password",
					},
				},
				Domains: []settings.Domain{
					{DomainName: "he-domain.com", SubDomains: []string{"www"}}, // Uses global HE
					{DomainName: "duck-domain.org", SubDomains: []string{"www"}, Provider: "DuckDNS"},
					{DomainName: "google-domain.com", SubDomains: []string{"www"}, Provider: "Google"},
					{DomainName: "ali-domain.com", SubDomains: []string{"www"}, Provider: "AliDNS"},
				},
			}

			providers, err := GetProviders(config)
			if err != nil {
				t.Fatalf("Failed to get providers: %v", err)
			}

			if len(providers) != 4 {
				t.Errorf("Expected 4 providers, got %d", len(providers))
			}

			expectedProviders := []string{"HE", "DuckDNS", "Google", "AliDNS"}
			for _, providerName := range expectedProviders {
				if _, exists := providers[providerName]; !exists {
					t.Errorf("%s provider not found", providerName)
				}
			}
		})
	})
}

func TestGetProviderForDomain(t *testing.T) {
	// Test single provider mode domain resolution
	t.Run("SingleProviderDomainResolution", func(t *testing.T) {
		config := &settings.Settings{
			Provider:   "DNSPod",
			LoginToken: "test-token",
		}

		providers, err := GetProviders(config)
		if err != nil {
			t.Fatalf("Failed to get providers: %v", err)
		}

		domain := &settings.Domain{
			DomainName: "example.com",
			SubDomains: []string{"www"},
		}

		provider, err := GetProviderForDomain(domain, providers, config)
		if err != nil {
			t.Fatalf("Failed to get provider for domain: %v", err)
		}

		if provider == nil {
			t.Error("Provider is nil")
		}
	})

	// Test multi-provider mode domain resolution
	t.Run("MultiProviderDomainResolution", func(t *testing.T) {
		config := &settings.Settings{
			Providers: map[string]*settings.ProviderConfig{
				"Cloudflare": {
					LoginToken: "cf-token",
				},
				"DNSPod": {
					LoginToken: "dnspod-token",
				},
			},
		}

		providers, err := GetProviders(config)
		if err != nil {
			t.Fatalf("Failed to get providers: %v", err)
		}

		// Test domain with specific provider
		domain1 := &settings.Domain{
			DomainName: "example.com",
			SubDomains: []string{"www"},
			Provider:   "Cloudflare",
		}

		provider1, err := GetProviderForDomain(domain1, providers, config)
		if err != nil {
			t.Fatalf("Failed to get provider for domain1: %v", err)
		}
		if provider1 == nil {
			t.Error("Provider1 is nil")
		}

		// Test domain with different provider
		domain2 := &settings.Domain{
			DomainName: "example.net",
			SubDomains: []string{"www"},
			Provider:   "DNSPod",
		}

		provider2, err := GetProviderForDomain(domain2, providers, config)
		if err != nil {
			t.Fatalf("Failed to get provider for domain2: %v", err)
		}
		if provider2 == nil {
			t.Error("Provider2 is nil")
		}
	})

	// Test mixed configuration domain resolution
	t.Run("MixedConfigurationDomainResolution", func(t *testing.T) {
		config := &settings.Settings{
			Provider:   "DNSPod",
			LoginToken: "global-token",
			Providers: map[string]*settings.ProviderConfig{
				"Cloudflare": {
					LoginToken: "cf-token",
				},
			},
		}

		providers, err := GetProviders(config)
		if err != nil {
			t.Fatalf("Failed to get providers: %v", err)
		}

		// Test domain using global provider (no provider specified)
		domain1 := &settings.Domain{
			DomainName: "old-domain.com",
			SubDomains: []string{"www"},
		}

		provider1, err := GetProviderForDomain(domain1, providers, config)
		if err != nil {
			t.Fatalf("Failed to get provider for global domain: %v", err)
		}
		if provider1 == nil {
			t.Error("Global provider is nil")
		}

		// Test domain using specific provider
		domain2 := &settings.Domain{
			DomainName: "new-domain.com",
			SubDomains: []string{"www"},
			Provider:   "Cloudflare",
		}

		provider2, err := GetProviderForDomain(domain2, providers, config)
		if err != nil {
			t.Fatalf("Failed to get provider for specific domain: %v", err)
		}
		if provider2 == nil {
			t.Error("Specific provider is nil")
		}
	})

	// Test error cases
	t.Run("ErrorCases", func(t *testing.T) {
		config := &settings.Settings{
			Providers: map[string]*settings.ProviderConfig{
				"Cloudflare": {
					LoginToken: "cf-token",
				},
			},
		}

		providers, err := GetProviders(config)
		if err != nil {
			t.Fatalf("Failed to get providers: %v", err)
		}

		// Test domain referencing non-existent provider
		domain := &settings.Domain{
			DomainName: "example.com",
			SubDomains: []string{"www"},
			Provider:   "NonExistentProvider",
		}

		_, err = GetProviderForDomain(domain, providers, config)
		if err == nil {
			t.Error("Expected error for non-existent provider, but got none")
		}
	})
}

func TestCreateProvider(t *testing.T) {
	// Test creating various provider types
	t.Run("CreateVariousProviders", func(t *testing.T) {
		testCases := []struct {
			name         string
			providerName string
			config       *settings.Settings
			shouldPass   bool
		}{
			{
				name:         "Cloudflare",
				providerName: "Cloudflare",
				config: &settings.Settings{
					Provider:   "Cloudflare",
					LoginToken: "test-token",
				},
				shouldPass: true,
			},
			{
				name:         "DNSPod",
				providerName: "DNSPod",
				config: &settings.Settings{
					Provider:   "DNSPod",
					LoginToken: "test-token",
				},
				shouldPass: true,
			},
			{
				name:         "DigitalOcean",
				providerName: "DigitalOcean",
				config: &settings.Settings{
					Provider:   "DigitalOcean",
					LoginToken: "test-token",
				},
				shouldPass: true,
			},
			{
				name:         "DuckDNS",
				providerName: "DuckDNS",
				config: &settings.Settings{
					Provider:   "DuckDNS",
					LoginToken: "test-token",
				},
				shouldPass: true,
			},
			{
				name:         "HE",
				providerName: "HE",
				config: &settings.Settings{
					Provider: "HE",
					Password: "test-password",
				},
				shouldPass: true,
			},
			{
				name:         "Google",
				providerName: "Google",
				config: &settings.Settings{
					Provider: "Google",
					Email:    "test@gmail.com",
					Password: "test-password",
				},
				shouldPass: true,
			},
			{
				name:         "OVH",
				providerName: "OVH",
				config: &settings.Settings{
					Provider:    "OVH",
					AppKey:      "test-key",
					AppSecret:   "test-secret",
					ConsumerKey: "test-consumer",
				},
				shouldPass: true,
			},
			{
				name:         "UnsupportedProvider",
				providerName: "UnsupportedProvider",
				config: &settings.Settings{
					Provider: "UnsupportedProvider",
				},
				shouldPass: false,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				provider, err := createProvider(tc.providerName, tc.config)
				if tc.shouldPass {
					if err != nil {
						t.Errorf("Expected success for %s, got error: %v", tc.name, err)
					}
					if provider == nil {
						t.Errorf("Expected non-nil provider for %s", tc.name)
					}
				} else {
					if err == nil {
						t.Errorf("Expected error for %s, but got none", tc.name)
					}
				}
			})
		}
	})
}
