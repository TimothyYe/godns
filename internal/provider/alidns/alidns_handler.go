package alidns

import (
	"fmt"

	"github.com/TimothyYe/godns/internal/settings"
	log "github.com/sirupsen/logrus"
)

// DNSProvider struct.
type DNSProvider struct {
	aliDNS *AliDNS
}

func (provider *DNSProvider) Init(conf *settings.Settings) {
	provider.aliDNS = NewAliDNS(
		conf.Email,
		conf.Password,
		conf.IPType)
}

func (provider *DNSProvider) UpdateIP(domainName, subdomainName, ip string) error {
	log.Infof("%s.%s - Start to update record IP...", subdomainName, domainName)
	records := provider.aliDNS.GetDomainRecords(domainName, subdomainName)
	if len(records) == 0 {
		log.Errorf("Cannot get subdomain [%s] from AliDNS.", subdomainName)
		return fmt.Errorf("cannot get subdomain [%s] from AliDNS", subdomainName)
	}

	if records[0].Value != ip {
		records[0].Value = ip
		if err := provider.aliDNS.UpdateDomainRecord(records[0]); err != nil {
			return fmt.Errorf("failed to update IP for subdomain: %s", subdomainName)
		}

		log.Infof("IP updated for subdomain: %s", subdomainName)
	} else {
		log.Debugf("IP not changed for subdomain: %s", subdomainName)
	}

	return nil
}
