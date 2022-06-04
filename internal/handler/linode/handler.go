package linode

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/linode/linodego"
	log "github.com/sirupsen/logrus"

	"github.com/TimothyYe/godns/internal/settings"
	"github.com/TimothyYe/godns/internal/utils"
	"github.com/TimothyYe/godns/pkg/notification"
)

type IDNSClient interface {
	UpdateDNSRecordIP(string, string, string) error
}

type Handler struct {
	Configuration *settings.Settings
	client        IDNSClient
	cachedIP      string
	notifyManager notification.INotificationManager
}

func (handler *Handler) SetConfiguration(conf *settings.Settings) {
	handler.Configuration = conf
	handler.client = createDNSClient(conf)
	handler.notifyManager = notification.GetNotificationManager(handler.Configuration)
}

func createDNSClient(conf *settings.Settings) IDNSClient {
	httpClient, err := CreateHTTPClient(conf)
	if err != nil {
		panic(err)
	}
	linodeAPIClient := linodego.NewClient(httpClient)
	linodeAPIClient.SetDebug(conf.DebugInfo)
	return CreateLinodeDNSClient(&linodeAPIClient)
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
	err = handler.updateDNS(domain, ip, handler.client)
	if err != nil {
		log.Error(err)
		return
	}
	handler.cachedIP = ip
	log.Debugf("Cached IP address: %s", ip)
}

func (handler *Handler) updateDNS(domain *settings.Domain, ip string, client IDNSClient) error {
	for _, subdomainName := range domain.SubDomains {
		err := client.UpdateDNSRecordIP(domain.DomainName, subdomainName, ip)
		if err != nil {
			return err
		}
		successMessage := fmt.Sprintf("%s.%s", subdomainName, domain.DomainName)
		handler.notifyManager.Send(successMessage, ip)
	}

	return nil
}
