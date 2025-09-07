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

// credentialAccessor defines an interface for accessing provider credentials.
type credentialAccessor interface {
	GetEmail() string
	GetPassword() string
	GetLoginToken() string
	GetAppKey() string
	GetAppSecret() string
	GetConsumerKey() string
}

// settingsAccessor adapts Settings to credentialAccessor interface.
type settingsAccessor struct {
	config *settings.Settings
}

func (s *settingsAccessor) GetEmail() string       { return s.config.Email }
func (s *settingsAccessor) GetPassword() string    { return s.config.Password }
func (s *settingsAccessor) GetLoginToken() string  { return s.config.LoginToken }
func (s *settingsAccessor) GetAppKey() string      { return s.config.AppKey }
func (s *settingsAccessor) GetAppSecret() string   { return s.config.AppSecret }
func (s *settingsAccessor) GetConsumerKey() string { return s.config.ConsumerKey }

// providerConfigAccessor adapts ProviderConfig to credentialAccessor interface.
type providerConfigAccessor struct {
	config *settings.ProviderConfig
}

func (p *providerConfigAccessor) GetEmail() string       { return p.config.Email }
func (p *providerConfigAccessor) GetPassword() string    { return p.config.Password }
func (p *providerConfigAccessor) GetLoginToken() string  { return p.config.LoginToken }
func (p *providerConfigAccessor) GetAppKey() string      { return p.config.AppKey }
func (p *providerConfigAccessor) GetAppSecret() string   { return p.config.AppSecret }
func (p *providerConfigAccessor) GetConsumerKey() string { return p.config.ConsumerKey }

// validateProviderCredentials validates provider credentials using the common interface.
func validateProviderCredentials(providerName string, accessor credentialAccessor) error {
	switch providerName {
	case DNSPOD:
		if accessor.GetPassword() == "" && accessor.GetLoginToken() == "" {
			return errors.New("password or login token cannot be empty")
		}
	case HE:
		if accessor.GetPassword() == "" {
			return errors.New("password cannot be empty")
		}
	case CLOUDFLARE:
		if accessor.GetLoginToken() == "" {
			if accessor.GetEmail() == "" {
				return errors.New("email cannot be empty")
			}
			if accessor.GetPassword() == "" {
				return errors.New("password cannot be empty")
			}
		}
	case ALIDNS:
		if accessor.GetEmail() == "" {
			return errors.New("email cannot be empty")
		}
		if accessor.GetPassword() == "" {
			return errors.New("password cannot be empty")
		}
	case DIGITALOCEAN:
		if accessor.GetLoginToken() == "" {
			return errors.New("login token cannot be empty")
		}
	case DUCK:
		if accessor.GetLoginToken() == "" {
			return errors.New("login token cannot be empty")
		}
	case DYNU:
		if accessor.GetPassword() == "" {
			return errors.New("password cannot be empty")
		}
	case DYNV6:
		if accessor.GetLoginToken() == "" {
			return errors.New("login token cannot be empty")
		}
	case GOOGLE:
		fallthrough
	case NOIP:
		if accessor.GetEmail() == "" {
			return errors.New("email cannot be empty")
		}
		if accessor.GetPassword() == "" {
			return errors.New("password cannot be empty")
		}
	case DREAMHOST:
		if accessor.GetLoginToken() == "" {
			return errors.New("login token cannot be empty")
		}
	case SCALEWAY:
		if accessor.GetLoginToken() == "" {
			return errors.New("login token cannot be empty")
		}
	case LINODE:
		if accessor.GetLoginToken() == "" {
			return errors.New("login token cannot be empty")
		}
	case STRATO:
		if accessor.GetPassword() == "" {
			return errors.New("password cannot be empty")
		}
	case LOOPIASE:
		if accessor.GetPassword() == "" {
			return errors.New("password cannot be empty")
		}
	case INFOMANIAK:
		if accessor.GetPassword() == "" {
			return errors.New("password cannot be empty")
		}
	case HETZNER:
		if accessor.GetLoginToken() == "" {
			return errors.New("login token cannot be empty")
		}
	case IONOS:
		if accessor.GetLoginToken() == "" {
			return errors.New("login token cannot be empty")
		}
	case OVH:
		if accessor.GetAppKey() == "" {
			return errors.New("app key cannot be empty")
		}
		if accessor.GetAppSecret() == "" {
			return errors.New("app secret cannot be empty")
		}
		if accessor.GetConsumerKey() == "" {
			return errors.New("consumer key cannot be empty")
		}
	case TRANSIP:
		if accessor.GetEmail() == "" {
			return errors.New("email cannot be empty")
		}
		if accessor.GetLoginToken() == "" {
			return errors.New("login token cannot be empty")
		}
	case PORKBUN:
		if accessor.GetLoginToken() == "" {
			return errors.New("API key cannot be empty")
		}
		if accessor.GetPassword() == "" {
			return errors.New("secret key cannot be empty")
		}
	default:
		return fmt.Errorf("'%s' is not a supported DNS provider", providerName)
	}

	return nil
}

// checkSingleProviderCredentials validates credentials for legacy single provider mode.
func checkSingleProviderCredentials(providerName string, config *settings.Settings) error {
	return validateProviderCredentials(providerName, &settingsAccessor{config})
}

// checkProviderCredentials validates credentials for a provider configuration.
func checkProviderCredentials(providerName string, providerConfig *settings.ProviderConfig) error {
	return validateProviderCredentials(providerName, &providerConfigAccessor{providerConfig})
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
