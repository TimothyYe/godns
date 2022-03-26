package linode

import (
	"runtime/debug"
	"time"

	"github.com/linode/linodego"
	log "github.com/sirupsen/logrus"

	"github.com/TimothyYe/godns/internal/settings"
	"github.com/TimothyYe/godns/internal/utils"
)

type DNSClient interface {
	GetDomainID(string) (int, error)
	GetDomainRecordID(int, string) (bool, int, error)
	CreateDomainRecord(int, string) (int, error)
	UpdateDomainRecord(int, int, string) error
}

type Handler struct {
	Configuration *settings.Settings
	client        DNSClient
	cachedIP      string
}

func (handler *Handler) SetConfiguration(conf *settings.Settings) {
	handler.Configuration = conf
	handler.client = createDNSClient(conf)
}

func createDNSClient(conf *settings.Settings) DNSClient {
	httpClient, err := CreateHTTPClient(conf)
	if err != nil {
		panic(err)
	}
	linodeApiClient := linodego.NewClient(httpClient)
	linodeApiClient.SetDebug(conf.DebugInfo)
	return CreateLinodeDNSClient(&linodeApiClient)
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
	err = updateDNS(domain, ip, handler.client)
	if err != nil {
		log.Error(err)
		return
	}
	handler.cachedIP = ip
	log.Debugf("Cached IP address: %s", ip)
}

func updateDNS(domain *settings.Domain, ip string, client DNSClient) error {
	domainID, err := client.GetDomainID(domain.DomainName)
	if err != nil {
		return err
	}

	for _, subdomainName := range domain.SubDomains {
		if subdomainName == utils.RootDomain {
			subdomainName = ""
		}
		recordExists, recordID, err := client.GetDomainRecordID(domainID, subdomainName)
		if err != nil {
			return err
		}
		if !recordExists {
			recordID, err = client.CreateDomainRecord(domainID, subdomainName)
		}

		err = client.UpdateDomainRecord(domainID, recordID, ip)
		if err != nil {
			return err
		}
	}

	return nil
}
