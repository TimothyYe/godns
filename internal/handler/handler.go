package handler

import (
	"errors"
	"fmt"
	"runtime/debug"
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
	cachedIP            string
}

func (handler *Handler) SetConfiguration(conf *settings.Settings) {
	handler.Configuration = conf
	handler.notificationManager = notification.GetNotificationManager(handler.Configuration)
}

func (handler *Handler) SetProvider(provider provider.IDNSProvider) {
	handler.dnsProvider = provider
}

func (handler *Handler) LoopUpdateIP(domain *settings.Domain, panicChan chan<- settings.Domain) error {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("Recovered in %v: %v", err, string(debug.Stack()))
			panicChan <- *domain
		}
	}()

	for {
		err := handler.UpdateIP(domain)
		if err != nil {
			log.WithError(err).Debug("Update IP failed during the DNS Update loop")
		}
		log.Debugf("DNS update loop finished, will run again in %d seconds", handler.Configuration.Interval)
		time.Sleep(time.Second * time.Duration(handler.Configuration.Interval))
	}
}

func (handler *Handler) UpdateIP(domain *settings.Domain) error {
	ip, err := utils.GetCurrentIP(handler.Configuration)
	if err != nil {
		if handler.Configuration.RunOnce {
			return fmt.Errorf("%v: fail to get current IP", err)
		}
		log.Error(err)
		return nil
	}

	if ip == handler.cachedIP {
		log.Debugf("IP (%s) matches cached IP (%s), skipping", ip, handler.cachedIP)
		return nil
	}

	err = handler.updateDNS(domain, ip)
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

			successMessage := fmt.Sprintf("%s.%s", subdomainName, domain.DomainName)
			handler.notificationManager.Send(successMessage, ip)

			// execute webhook when it is enabled
			if handler.Configuration.Webhook.Enabled {
				if err := lib.GetWebhook(handler.Configuration).Execute(hostname, ip); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
