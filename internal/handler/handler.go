package handler

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/TimothyYe/godns/internal/provider"

	log "github.com/sirupsen/logrus"

	"github.com/TimothyYe/godns/internal/settings"
	"github.com/TimothyYe/godns/internal/utils"
	"github.com/TimothyYe/godns/pkg/lib"
	"github.com/TimothyYe/godns/pkg/notification"
)

var (
	errEmptyResult = errors.New("empty result")
	errEmptyDomain = errors.New("NXDOMAIN")
)

type Handler struct {
	ctx                 context.Context
	Configuration       *settings.Settings
	dnsProvider         provider.IDNSProvider            // Legacy single provider
	dnsProviders        map[string]provider.IDNSProvider // Multi-provider support
	notificationManager notification.INotificationManager
	ipManager           *lib.IPHelper
	cachedIP            string
}

func (handler *Handler) SetContext(ctx context.Context) {
	handler.ctx = ctx
}

func (handler *Handler) SetConfiguration(conf *settings.Settings) {
	handler.Configuration = conf
	handler.notificationManager = notification.GetNotificationManager(handler.Configuration)
	handler.ipManager = lib.GetIPHelperInstance(handler.Configuration)
}

func (handler *Handler) Init() {
	handler.ipManager.UpdateConfiguration(handler.Configuration)
}

func (handler *Handler) SetProvider(provider provider.IDNSProvider) {
	handler.dnsProvider = provider
}

func (handler *Handler) SetProviders(providers map[string]provider.IDNSProvider) {
	handler.dnsProviders = providers
}

func (handler *Handler) LoopUpdateIP(ctx context.Context, domain *settings.Domain) error {
	ticker := time.NewTicker(time.Second * time.Duration(handler.Configuration.Interval))

	// run once at the beginning
	err := handler.UpdateIP(domain)
	if err != nil {
		log.WithError(err).Debug("Update IP failed during the DNS Update loop")
	}
	log.Debugf("DNS update loop finished, will run again in %d seconds", handler.Configuration.Interval)

	for {
		select {
		case <-ticker.C:
			err := handler.UpdateIP(domain)
			if err != nil {
				log.WithError(err).Debug("Update IP failed during the DNS Update loop")
			}
			log.Debugf("DNS update loop finished, will run again in %d seconds", handler.Configuration.Interval)
		case <-ctx.Done():
			log.Info("DNS update loop cancelled")
			ticker.Stop()
			return nil
		}
	}
}

func (handler *Handler) UpdateIP(domain *settings.Domain) error {
	ip := handler.ipManager.GetCurrentIP()
	if ip == "" {
		if handler.Configuration.RunOnce {
			return fmt.Errorf("fail to get current IP")
		}
		return nil
	}

	if ip == handler.cachedIP {
		log.Debugf("IP (%s) matches cached IP (%s), skipping", ip, handler.cachedIP)
		return nil
	}

	err := handler.updateDNS(domain, ip)
	if err != nil {
		if handler.Configuration.RunOnce {
			return fmt.Errorf("%v: fail to update DNS", err)
		}
		log.Error(err)
		return nil
	}
	handler.cachedIP = ip
	log.Debugf("Cached IP address: %s", ip)
	return nil
}

func (handler *Handler) updateDNS(domain *settings.Domain, ip string) error {
	var updatedDomains []string

	// Get the appropriate provider for this domain
	domainProvider, err := handler.getProviderForDomain(domain)
	if err != nil {
		return fmt.Errorf("failed to get provider for domain %s: %w", domain.DomainName, err)
	}

	for _, subdomainName := range domain.SubDomains {
		var hostname string

		if subdomainName != utils.RootDomain {
			hostname = subdomainName + "." + domain.DomainName
		} else {
			hostname = domain.DomainName
		}

		lastIP, err := utils.ResolveDNS(hostname, handler.Configuration.Resolver, handler.Configuration.IPType)
		if err != nil && (errors.Is(err, errEmptyResult) || errors.Is(err, errEmptyDomain)) {
			log.Errorf("Failed to resolve DNS for domain: %s, error: %s", hostname, err)
			continue
		}

		// check against the current known IP, if no change, skip update
		if ip == lastIP {
			log.Infof("IP is the same as cached one (%s). Skip update.", ip)
		} else {

			if err := domainProvider.UpdateIP(domain.DomainName, subdomainName, ip); err != nil {
				return err
			}

			updatedDomains = append(updatedDomains, subdomainName)

			// execute webhook when it is enabled
			if handler.Configuration.Webhook.Enabled {
				if err := lib.GetWebhook(handler.Configuration).Execute(hostname, ip); err != nil {
					return err
				}
			}
		}
	}

	if len(updatedDomains) > 0 {
		providerName := handler.Configuration.GetDomainProvider(domain)
		successMessage := fmt.Sprintf("[ %s ] of %s (via %s)", strings.Join(updatedDomains, ", "), domain.DomainName, providerName)
		handler.notificationManager.Send(successMessage, ip)
	}

	return nil
}

// getProviderForDomain returns the appropriate provider for a given domain.
func (handler *Handler) getProviderForDomain(domain *settings.Domain) (provider.IDNSProvider, error) {
	// Multi-provider mode
	if handler.Configuration.IsMultiProvider() {
		providerName := handler.Configuration.GetDomainProvider(domain)
		domainProvider, exists := handler.dnsProviders[providerName]
		if !exists {
			return nil, fmt.Errorf("provider '%s' not found for domain %s", providerName, domain.DomainName)
		}
		return domainProvider, nil
	}

	// Legacy single provider mode
	if handler.dnsProvider == nil {
		return nil, fmt.Errorf("no DNS provider configured")
	}
	return handler.dnsProvider, nil
}
