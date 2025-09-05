package provider

import (
	"fmt"

	"github.com/TimothyYe/godns/internal/provider/alidns"
	"github.com/TimothyYe/godns/internal/provider/cloudflare"
	"github.com/TimothyYe/godns/internal/provider/digitalocean"
	"github.com/TimothyYe/godns/internal/provider/dnspod"
	"github.com/TimothyYe/godns/internal/provider/dreamhost"
	"github.com/TimothyYe/godns/internal/provider/duck"
	"github.com/TimothyYe/godns/internal/provider/dynu"
	"github.com/TimothyYe/godns/internal/provider/dynv6"
	"github.com/TimothyYe/godns/internal/provider/google"
	"github.com/TimothyYe/godns/internal/provider/he"
	"github.com/TimothyYe/godns/internal/provider/hetzner"
	"github.com/TimothyYe/godns/internal/provider/infomaniak"
	"github.com/TimothyYe/godns/internal/provider/ionos"
	"github.com/TimothyYe/godns/internal/provider/linode"
	"github.com/TimothyYe/godns/internal/provider/loopiase"
	"github.com/TimothyYe/godns/internal/provider/noip"
	"github.com/TimothyYe/godns/internal/provider/ovh"
	"github.com/TimothyYe/godns/internal/provider/scaleway"
	"github.com/TimothyYe/godns/internal/provider/strato"
	"github.com/TimothyYe/godns/internal/provider/transip"
	"github.com/TimothyYe/godns/internal/settings"
	"github.com/TimothyYe/godns/internal/utils"
)

func GetProvider(conf *settings.Settings) (IDNSProvider, error) {
	return createProvider(conf.Provider, conf)
}

// GetProviders returns a map of all configured providers for multi-provider support
func GetProviders(conf *settings.Settings) (map[string]IDNSProvider, error) {
	providers := make(map[string]IDNSProvider)
	
	// Handle legacy single provider mode
	if !conf.IsMultiProvider() {
		if conf.Provider == "" {
			return nil, fmt.Errorf("no provider configured")
		}
		provider, err := createProvider(conf.Provider, conf)
		if err != nil {
			return nil, err
		}
		providers[conf.Provider] = provider
		return providers, nil
	}
	
	// Handle multi-provider mode
	for providerName, providerConfig := range conf.Providers {
		// Create a temporary settings object with provider-specific config
		tempSettings := *conf
		tempSettings.Provider = providerName
		tempSettings.Email = providerConfig.Email
		tempSettings.Password = providerConfig.Password
		tempSettings.PasswordFile = providerConfig.PasswordFile
		tempSettings.LoginToken = providerConfig.LoginToken
		tempSettings.LoginTokenFile = providerConfig.LoginTokenFile
		tempSettings.AppKey = providerConfig.AppKey
		tempSettings.AppSecret = providerConfig.AppSecret
		tempSettings.ConsumerKey = providerConfig.ConsumerKey
		
		provider, err := createProvider(providerName, &tempSettings)
		if err != nil {
			return nil, fmt.Errorf("failed to create provider %s: %w", providerName, err)
		}
		
		providers[providerName] = provider
	}
	
	return providers, nil
}

// GetProviderForDomain returns the appropriate provider for a given domain
func GetProviderForDomain(domain *settings.Domain, providers map[string]IDNSProvider, conf *settings.Settings) (IDNSProvider, error) {
	providerName := conf.GetDomainProvider(domain)
	
	provider, exists := providers[providerName]
	if !exists {
		return nil, fmt.Errorf("provider '%s' not found for domain %s", providerName, domain.DomainName)
	}
	
	return provider, nil
}

func createProvider(providerName string, conf *settings.Settings) (IDNSProvider, error) {
	var provider IDNSProvider

	switch providerName {
	case utils.CLOUDFLARE:
		provider = &cloudflare.DNSProvider{}
	case utils.DIGITALOCEAN:
		provider = &digitalocean.DNSProvider{}
	case utils.DNSPOD:
		provider = &dnspod.DNSProvider{}
	case utils.DREAMHOST:
		provider = &dreamhost.DNSProvider{}
	case utils.HE:
		provider = &he.DNSProvider{}
	case utils.ALIDNS:
		provider = &alidns.DNSProvider{}
	case utils.GOOGLE:
		provider = &google.DNSProvider{}
	case utils.DUCK:
		provider = &duck.DNSProvider{}
	case utils.NOIP:
		provider = &noip.DNSProvider{}
	case utils.SCALEWAY:
		provider = &scaleway.DNSProvider{}
	case utils.DYNV6:
		provider = &dynv6.DNSProvider{}
	case utils.LINODE:
		provider = &linode.DNSProvider{}
	case utils.STRATO:
		provider = &strato.DNSProvider{}
	case utils.LOOPIASE:
		provider = &loopiase.DNSProvider{}
	case utils.INFOMANIAK:
		provider = &infomaniak.DNSProvider{}
	case utils.HETZNER:
		provider = &hetzner.DNSProvider{}
	case utils.OVH:
		provider = &ovh.DNSProvider{}
	case utils.DYNU:
		provider = &dynu.DNSProvider{}
	case utils.IONOS:
		provider = &ionos.DNSProvider{}
	case utils.TRANSIP:
		provider = &transip.DNSProvider{}
	default:
		return nil, fmt.Errorf("unknown provider '%s'", providerName)
	}

	provider.Init(conf)
	return provider, nil
}
