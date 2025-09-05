package utils

import (
	"testing"

	"github.com/TimothyYe/godns/internal/settings"
)

func TestCheckSettings(t *testing.T) {
	// Test legacy single-provider mode
	t.Run("LegacySingleProvider", func(t *testing.T) {
		// Empty settings should fail
		settingError := &settings.Settings{}
		if err := CheckSettings(settingError); err == nil {
			t.Error("empty setting should return error")
		}

		// Valid DNSPod with login token should pass
		settingDNSPod := &settings.Settings{
			Provider:   "DNSPod",
			LoginToken: "test-token",
			Domains: []settings.Domain{
				{DomainName: "example.com", SubDomains: []string{"www"}},
			},
		}
		if err := CheckSettings(settingDNSPod); err != nil {
			t.Errorf("valid DNSPod setting should pass, got error: %v", err)
		}

		// DNSPod without credentials should fail
		settingDNSPodInvalid := &settings.Settings{
			Provider: "DNSPod",
			Domains: []settings.Domain{
				{DomainName: "example.com", SubDomains: []string{"www"}},
			},
		}
		if err := CheckSettings(settingDNSPodInvalid); err == nil {
			t.Error("DNSPod setting without credentials should fail")
		}

		// HE without password should fail
		settingHE := &settings.Settings{
			Provider: "HE",
			Password: "",
			Domains: []settings.Domain{
				{DomainName: "example.com", SubDomains: []string{"www"}},
			},
		}
		if err := CheckSettings(settingHE); err == nil {
			t.Error("HE setting without password should fail")
		}

		// Cloudflare with email and password should pass
		settingCloudflare := &settings.Settings{
			Provider: "Cloudflare",
			Email:    "test@example.com",
			Password: "test-password",
			Domains: []settings.Domain{
				{DomainName: "example.com", SubDomains: []string{"www"}},
			},
		}
		if err := CheckSettings(settingCloudflare); err != nil {
			t.Errorf("valid Cloudflare setting should pass, got error: %v", err)
		}

		// Cloudflare with API token should pass
		settingCloudflareToken := &settings.Settings{
			Provider:   "Cloudflare",
			LoginToken: "test-token",
			Domains: []settings.Domain{
				{DomainName: "example.com", SubDomains: []string{"www"}},
			},
		}
		if err := CheckSettings(settingCloudflareToken); err != nil {
			t.Errorf("Cloudflare setting with token should pass, got error: %v", err)
		}

		// OVH with all required keys should pass
		settingOVH := &settings.Settings{
			Provider:    "OVH",
			AppKey:      "test-app-key",
			AppSecret:   "test-app-secret",
			ConsumerKey: "test-consumer-key",
			Domains: []settings.Domain{
				{DomainName: "example.com", SubDomains: []string{"www"}},
			},
		}
		if err := CheckSettings(settingOVH); err != nil {
			t.Errorf("valid OVH setting should pass, got error: %v", err)
		}
	})

	// Test multi-provider mode
	t.Run("MultiProvider", func(t *testing.T) {
		// Valid multi-provider configuration
		settingMulti := &settings.Settings{
			Providers: map[string]*settings.ProviderConfig{
				"Cloudflare": {
					Email:    "test@example.com",
					Password: "cf-token",
				},
				"DNSPod": {
					LoginToken: "dnspod-token",
				},
			},
			Domains: []settings.Domain{
				{DomainName: "example.com", SubDomains: []string{"www"}, Provider: "Cloudflare"},
				{DomainName: "example.net", SubDomains: []string{"www"}, Provider: "DNSPod"},
			},
		}
		if err := CheckSettings(settingMulti); err != nil {
			t.Errorf("valid multi-provider setting should pass, got error: %v", err)
		}

		// Multi-provider with empty providers map should fail
		settingMultiEmpty := &settings.Settings{
			Providers: map[string]*settings.ProviderConfig{},
			Domains: []settings.Domain{
				{DomainName: "example.com", SubDomains: []string{"www"}},
			},
		}
		if err := CheckSettings(settingMultiEmpty); err == nil {
			t.Error("multi-provider with empty providers should fail")
		}

		// Multi-provider with invalid provider credentials should fail
		settingMultiInvalid := &settings.Settings{
			Providers: map[string]*settings.ProviderConfig{
				"Cloudflare": {
					// Missing credentials
				},
			},
			Domains: []settings.Domain{
				{DomainName: "example.com", SubDomains: []string{"www"}, Provider: "Cloudflare"},
			},
		}
		if err := CheckSettings(settingMultiInvalid); err == nil {
			t.Error("multi-provider with invalid credentials should fail")
		}

		// Domain references non-existent provider should fail
		settingMissingProvider := &settings.Settings{
			Providers: map[string]*settings.ProviderConfig{
				"Cloudflare": {
					LoginToken: "cf-token",
				},
			},
			Domains: []settings.Domain{
				{DomainName: "example.com", SubDomains: []string{"www"}, Provider: "nonexistent"},
			},
		}
		if err := CheckSettings(settingMissingProvider); err == nil {
			t.Error("domain referencing non-existent provider should fail")
		}
	})

	// Test mixed configuration (global provider + per-domain providers)
	t.Run("MixedConfiguration", func(t *testing.T) {
		// Valid mixed configuration
		settingMixed := &settings.Settings{
			Provider:   "DNSPod",
			LoginToken: "global-dnspod-token",
			Providers: map[string]*settings.ProviderConfig{
				"Cloudflare": {
					LoginToken: "cf-token",
				},
			},
			Domains: []settings.Domain{
				{DomainName: "old-domain.com", SubDomains: []string{"www"}},                         // Uses global DNSPod
				{DomainName: "new-domain.com", SubDomains: []string{"www"}, Provider: "Cloudflare"}, // Uses specific provider
			},
		}
		if err := CheckSettings(settingMixed); err != nil {
			t.Errorf("valid mixed configuration should pass, got error: %v", err)
		}

		// Mixed configuration with domain having no provider and no global provider should fail
		settingMixedNoGlobal := &settings.Settings{
			Providers: map[string]*settings.ProviderConfig{
				"Cloudflare": {
					LoginToken: "cf-token",
				},
			},
			Domains: []settings.Domain{
				{DomainName: "no-provider.com", SubDomains: []string{"www"}}, // No provider specified and no global
			},
		}
		if err := CheckSettings(settingMixedNoGlobal); err == nil {
			t.Error("domain without provider and no global provider should fail")
		}
	})

	// Test domain validation
	t.Run("DomainValidation", func(t *testing.T) {
		// Empty domain name should fail
		settingEmptyDomain := &settings.Settings{
			Provider:   "DNSPod",
			LoginToken: "test-token",
			Domains: []settings.Domain{
				{DomainName: "", SubDomains: []string{"www"}},
			},
		}
		if err := CheckSettings(settingEmptyDomain); err == nil {
			t.Error("empty domain name should fail")
		}

		// Empty subdomain should fail
		settingEmptySubdomain := &settings.Settings{
			Provider:   "DNSPod",
			LoginToken: "test-token",
			Domains: []settings.Domain{
				{DomainName: "example.com", SubDomains: []string{""}},
			},
		}
		if err := CheckSettings(settingEmptySubdomain); err == nil {
			t.Error("empty subdomain should fail")
		}

		// No domains should fail
		settingNoDomains := &settings.Settings{
			Provider:   "DNSPod",
			LoginToken: "test-token",
			Domains:    []settings.Domain{},
		}
		if err := CheckSettings(settingNoDomains); err == nil {
			t.Error("no domains should fail")
		}
	})

	// Test provider-specific validations
	t.Run("ProviderSpecificValidation", func(t *testing.T) {
		// Test various providers with their specific requirements
		testCases := []struct {
			name        string
			config      *settings.ProviderConfig
			shouldPass  bool
			description string
		}{
			{
				name: "DNSPod",
				config: &settings.ProviderConfig{
					LoginToken: "test-token",
				},
				shouldPass:  true,
				description: "DNSPod with login token",
			},
			{
				name: "DNSPod",
				config: &settings.ProviderConfig{
					Password: "test-password",
				},
				shouldPass:  true,
				description: "DNSPod with password",
			},
			{
				name:   "DNSPod",
				config: &settings.ProviderConfig{
					// No credentials
				},
				shouldPass:  false,
				description: "DNSPod without credentials",
			},
			{
				name: "HE",
				config: &settings.ProviderConfig{
					Password: "test-password",
				},
				shouldPass:  true,
				description: "HE with password",
			},
			{
				name:   "HE",
				config: &settings.ProviderConfig{
					// No password
				},
				shouldPass:  false,
				description: "HE without password",
			},
			{
				name: "DigitalOcean",
				config: &settings.ProviderConfig{
					LoginToken: "do-token",
				},
				shouldPass:  true,
				description: "DigitalOcean with token",
			},
			{
				name: "AliDNS",
				config: &settings.ProviderConfig{
					Email:    "test@example.com",
					Password: "test-secret",
				},
				shouldPass:  true,
				description: "AliDNS with email and password",
			},
			{
				name: "AliDNS",
				config: &settings.ProviderConfig{
					Email: "test@example.com",
					// Missing password
				},
				shouldPass:  false,
				description: "AliDNS missing password",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.description, func(t *testing.T) {
				setting := &settings.Settings{
					Providers: map[string]*settings.ProviderConfig{
						tc.name: tc.config,
					},
					Domains: []settings.Domain{
						{DomainName: "example.com", SubDomains: []string{"www"}, Provider: tc.name},
					},
				}

				err := CheckSettings(setting)
				if tc.shouldPass && err != nil {
					t.Errorf("%s should pass but got error: %v", tc.description, err)
				}
				if !tc.shouldPass && err == nil {
					t.Errorf("%s should fail but passed", tc.description)
				}
			})
		}
	})
}
