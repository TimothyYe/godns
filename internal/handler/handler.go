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
	Configuration       *settings.Settings
	dnsProvider         provider.IDNSProvider
	notificationManager notification.INotificationManager
	ipManager           *lib.IPHelper
	cachedIP            string
}

func (handler *Handler) SetConfiguration(conf *settings.Settings) {
	handler.Configuration = conf
	handler.notificationManager = notification.GetNotificationManager(handler.Configuration)
	handler.ipManager = lib.NewIPHelper(handler.Configuration)
}

func (handler *Handler) SetProvider(provider provider.IDNSProvider) {
	handler.dnsProvider = provider
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

		//check against the current known IP, if no change, skip update
		if ip == lastIP {
			log.Infof("IP is the same as cached one (%s). Skip update.", ip)
		} else {

			if err := handler.dnsProvider.UpdateIP(domain.DomainName, subdomainName, ip); err != nil {
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
		successMessage := fmt.Sprintf("[ %s ] of %s", strings.Join(updatedDomains, ", "), domain.DomainName)
		handler.notificationManager.Send(successMessage, ip)
	}

	return nil
}
