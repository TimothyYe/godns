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

func (provider *DNSProvider) UpdateIP(domainName string, subdomainName string, ip string) error {
	log.Infof("%s.%s - Start to update record IP...", subdomainName, domainName)
	records := provider.aliDNS.GetDomainRecords(domainName, subdomainName)
	if len(records) == 0 {
		log.Errorf("Cannot get subdomain %s from AliDNS.", subdomainName)
		return fmt.Errorf("cannot get subdomain %s from AliDNS", subdomainName)
	}

	records[0].Value = ip
	if err := provider.aliDNS.UpdateDomainRecord(records[0]); err != nil {
		log.Errorf("Failed to update IP for subdomain: %s", subdomainName)
		return fmt.Errorf("failed to update IP for subdomain: %s", subdomainName)
	} else {
		log.Infof("IP updated for subdomain:%s", subdomainName)
	}

	return nil
}
