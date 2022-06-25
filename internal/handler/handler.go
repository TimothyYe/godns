package handler

import (
	"fmt"
	"runtime/debug"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/TimothyYe/godns/internal/settings"
	"github.com/TimothyYe/godns/internal/utils"
	"github.com/TimothyYe/godns/pkg/notification"
)

type Handler struct {
	Configuration       *settings.Settings
	dnsProvider         IDNSProvider
	notificationManager notification.INotificationManager
	cachedIP            string
}

func (handler *Handler) SetConfiguration(conf *settings.Settings) {
	handler.Configuration = conf
	handler.notificationManager = notification.GetNotificationManager(handler.Configuration)
}

func (handler *Handler) SetProvider(provider IDNSProvider) {
	handler.dnsProvider = provider
}

func (handler *Handler) DomainLoop(domain *settings.Domain, panicChan chan<- settings.Domain, runOnce bool) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("Recovered in %v: %v", err, string(debug.Stack()))
			panicChan <- *domain
		}
	}()

	for while := true; while; while = !runOnce {
		handler.domainLoop(domain)
		log.Debugf("DNS update loop finished, will run again in %d seconds", handler.Configuration.Interval)
		time.Sleep(time.Second * time.Duration(handler.Configuration.Interval))
	}
}

func (handler *Handler) domainLoop(domain *settings.Domain) {
	ip, err := utils.GetCurrentIP(handler.Configuration)
	if err != nil {
		log.Error(err)
		return
	}
	if ip == handler.cachedIP {
		log.Debugf("IP (%s) matches cached IP (%s), skipping", ip, handler.cachedIP)
		return
	}
	err = handler.updateDNS(domain, ip)
	if err != nil {
		log.Error(err)
		return
	}
	handler.cachedIP = ip
	log.Debugf("Cached IP address: %s", ip)
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
		if err != nil {
			log.Error(err)
			continue
		}

		//check against the current known IP, if no change, skip update
		if ip == lastIP {
			log.Infof("IP is the same as cached one (%s). Skip update.", ip)
		} else {
			err := handler.dnsProvider.UpdateIP(domain.DomainName, subdomainName, ip)
			if err != nil {
				return err
			}
			successMessage := fmt.Sprintf("%s.%s", subdomainName, domain.DomainName)
			handler.notificationManager.Send(successMessage, ip)
		}
	}

	return nil
}
