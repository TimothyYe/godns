package utils

import (
	"errors"
	"fmt"

	"github.com/TimothyYe/godns/internal/settings"
)

// CheckSettings check the format of settings.
func CheckSettings(config *settings.Settings) error {
	switch config.Provider {
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
	case DUCK:
		if config.LoginToken == "" {
			return errors.New("login token cannot be empty")
		}
	case DYNV6:
		if config.LoginToken == "" {
			return errors.New("login token cannot be empty")
		}
	case GOOGLE:
		any := false
		if config.Email == "" {
			for _, domain := range config.Domains {
				for _, sub_domain := range domain.DelegatedSubDomains {
					if len(sub_domain.Email) == 0 {
						name := fmt.Sprint("failed to find delegated sub domain email: %s:%s", domain.DomainName, sub_domain.DomainName)
						return errors.New(name)
					} else {
						any = true
					}

				}
			}
			if !any {
				return errors.New("email cannot be empty")
			}
		}
		if config.Password == "" && any == false {
			return errors.New("password cannot be empty")
		}
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

	default:
		message := fmt.Sprintf("'%s' is not a supported DNS provider", config.Provider)
		return errors.New(message)

	}

	return checkDomains(config)
}

func checkDomains(config *settings.Settings) error {
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
