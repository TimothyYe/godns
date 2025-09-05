package utils

import (
	"errors"
	"fmt"

	"github.com/TimothyYe/godns/internal/settings"
)

// CheckSettings check the format of settings.
func CheckSettings(config *settings.Settings) error {
	// Check if it's multi-provider mode
	if config.IsMultiProvider() {
		return checkMultiProviderSettings(config)
	}

	// Legacy single provider mode validation
	if config.Provider == "" {
		return errors.New("provider cannot be empty in single-provider mode")
	}

	if err := checkSingleProviderCredentials(config.Provider, config); err != nil {
		return err
	}

	return checkDomains(config)
}

// checkMultiProviderSettings validates multi-provider configuration.
func checkMultiProviderSettings(config *settings.Settings) error {
	if len(config.Providers) == 0 {
		return errors.New("providers cannot be empty in multi-provider mode")
	}

	// Validate each provider configuration
	for providerName, providerConfig := range config.Providers {
		if err := checkProviderCredentials(providerName, providerConfig); err != nil {
			return fmt.Errorf("provider '%s': %w", providerName, err)
		}
	}

	// Validate that all domain providers are configured
	return checkDomainsWithProviders(config)
}

// checkSingleProviderCredentials validates credentials for legacy single provider mode.
func checkSingleProviderCredentials(providerName string, config *settings.Settings) error {
	switch providerName {
	case DNSPOD:
		if config.Password == "" && config.LoginToken == "" {
			return errors.New("password or login token cannot be empty")
		}
	case HE:
		if config.Password == "" {
			return errors.New("password cannot be empty")
		}
	case CLOUDFLARE:
		if config.LoginToken == "" {
			if config.Email == "" {
				return errors.New("email cannot be empty")
			}
			if config.Password == "" {
				return errors.New("password cannot be empty")
			}
		}
	case ALIDNS:
		if config.Email == "" {
			return errors.New("email cannot be empty")
		}
		if config.Password == "" {
			return errors.New("password cannot be empty")
		}
	case DIGITALOCEAN:
		if config.LoginToken == "" {
			return errors.New("login token cannot be empty")
		}
	case DUCK:
		if config.LoginToken == "" {
			return errors.New("login token cannot be empty")
		}
	case DYNU:
		if config.Password == "" {
			return errors.New("password cannot be empty")
		}
	case DYNV6:
		if config.LoginToken == "" {
			return errors.New("login token cannot be empty")
		}
	case GOOGLE:
		fallthrough
	case NOIP:
		if config.Email == "" {
			return errors.New("email cannot be empty")
		}
		if config.Password == "" {
			return errors.New("password cannot be empty")
		}
	case DREAMHOST:
		if config.LoginToken == "" {
			return errors.New("login token cannot be empty")
		}
	case SCALEWAY:
		if config.LoginToken == "" {
			return errors.New("login token cannot be empty")
		}
	case LINODE:
		if config.LoginToken == "" {
			return errors.New("login token cannot be empty")
		}
	case STRATO:
		if config.Password == "" {
			return errors.New("password cannot be empty")
		}
	case LOOPIASE:
		if config.Password == "" {
			return errors.New("password cannot be empty")
		}
	case INFOMANIAK:
		if config.Password == "" {
			return errors.New("password cannot be empty")
		}
	case HETZNER:
		if config.LoginToken == "" {
			return errors.New("login token cannot be empty")
		}
	case IONOS:
		if config.LoginToken == "" {
			return errors.New("login token cannot be empty")
		}
	case OVH:
		if config.AppKey == "" {
			return errors.New("app key cannot be empty")
		}
		if config.AppSecret == "" {
			return errors.New("app secret cannot be empty")
		}
		if config.ConsumerKey == "" {
			return errors.New("consumer key cannot be empty")
		}
	case TRANSIP:
		if config.Email == "" {
			return errors.New("email cannot be empty")
		}
		if config.LoginToken == "" {
			return errors.New("login token cannot be empty")
		}
	default:
		return fmt.Errorf("'%s' is not a supported DNS provider", providerName)
	}

	return nil
}

// checkProviderCredentials validates credentials for a provider configuration.
func checkProviderCredentials(providerName string, providerConfig *settings.ProviderConfig) error {
	switch providerName {
	case DNSPOD:
		if providerConfig.Password == "" && providerConfig.LoginToken == "" {
			return errors.New("password or login token cannot be empty")
		}
	case HE:
		if providerConfig.Password == "" {
			return errors.New("password cannot be empty")
		}
	case CLOUDFLARE:
		if providerConfig.LoginToken == "" {
			if providerConfig.Email == "" {
				return errors.New("email cannot be empty")
			}
			if providerConfig.Password == "" {
				return errors.New("password cannot be empty")
			}
		}
	case ALIDNS:
		if providerConfig.Email == "" {
			return errors.New("email cannot be empty")
		}
		if providerConfig.Password == "" {
			return errors.New("password cannot be empty")
		}
	case DIGITALOCEAN:
		if providerConfig.LoginToken == "" {
			return errors.New("login token cannot be empty")
		}
	case DUCK:
		if providerConfig.LoginToken == "" {
			return errors.New("login token cannot be empty")
		}
	case DYNU:
		if providerConfig.Password == "" {
			return errors.New("password cannot be empty")
		}
	case DYNV6:
		if providerConfig.LoginToken == "" {
			return errors.New("login token cannot be empty")
		}
	case GOOGLE:
		fallthrough
	case NOIP:
		if providerConfig.Email == "" {
			return errors.New("email cannot be empty")
		}
		if providerConfig.Password == "" {
			return errors.New("password cannot be empty")
		}
	case DREAMHOST:
		if providerConfig.LoginToken == "" {
			return errors.New("login token cannot be empty")
		}
	case SCALEWAY:
		if providerConfig.LoginToken == "" {
			return errors.New("login token cannot be empty")
		}
	case LINODE:
		if providerConfig.LoginToken == "" {
			return errors.New("login token cannot be empty")
		}
	case STRATO:
		if providerConfig.Password == "" {
			return errors.New("password cannot be empty")
		}
	case LOOPIASE:
		if providerConfig.Password == "" {
			return errors.New("password cannot be empty")
		}
	case INFOMANIAK:
		if providerConfig.Password == "" {
			return errors.New("password cannot be empty")
		}
	case HETZNER:
		if providerConfig.LoginToken == "" {
			return errors.New("login token cannot be empty")
		}
	case IONOS:
		if providerConfig.LoginToken == "" {
			return errors.New("login token cannot be empty")
		}
	case OVH:
		if providerConfig.AppKey == "" {
			return errors.New("app key cannot be empty")
		}
		if providerConfig.AppSecret == "" {
			return errors.New("app secret cannot be empty")
		}
		if providerConfig.ConsumerKey == "" {
			return errors.New("consumer key cannot be empty")
		}
	case TRANSIP:
		if providerConfig.Email == "" {
			return errors.New("email cannot be empty")
		}
		if providerConfig.LoginToken == "" {
			return errors.New("login token cannot be empty")
		}
	default:
		return fmt.Errorf("'%s' is not a supported DNS provider", providerName)
	}

	return nil
}

// checkDomainsWithProviders validates domains in multi-provider mode.
func checkDomainsWithProviders(config *settings.Settings) error {
	if len(config.Domains) == 0 {
		return errors.New("at least one domain must be configured")
	}

	for _, d := range config.Domains {
		if d.DomainName == "" {
			return errors.New("domain name should not be empty")
		}

		// Validate subdomains
		for _, sd := range d.SubDomains {
			if sd == "" {
				return errors.New("subdomain should not be empty")
			}
		}

		// Get the provider for this domain (either domain-specific or global fallback)
		providerName := config.GetDomainProvider(&d)
		if providerName == "" {
			return fmt.Errorf("no provider configured for domain '%s'", d.DomainName)
		}

		// Check if the provider is configured
		if d.Provider != "" {
			// Domain has specific provider - check if it's configured in providers
			if _, exists := config.Providers[d.Provider]; !exists {
				return fmt.Errorf("domain '%s' references provider '%s' which is not configured in providers section", d.DomainName, d.Provider)
			}
		} else if config.Provider == "" {
			// No domain-specific provider and no global provider
			return fmt.Errorf("domain '%s' has no provider specified and no global provider is configured", d.DomainName)
		}
		// If domain doesn't specify provider but global provider exists, that's valid (mixed mode)
	}

	return nil
}

func checkDomains(config *settings.Settings) error {
	if len(config.Domains) == 0 {
		return errors.New("at least one domain must be configured")
	}

	for _, d := range config.Domains {
		if d.DomainName == "" {
			return errors.New("domain name should not be empty")
		}

		for _, sd := range d.SubDomains {
			if sd == "" {
				return errors.New("subdomain should not be empty")
			}
		}
	}

	return nil
}
