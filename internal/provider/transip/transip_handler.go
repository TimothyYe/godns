package transip

import (
	"strings"

	"github.com/TimothyYe/godns/internal/settings"
	"github.com/TimothyYe/godns/internal/utils"
	log "github.com/sirupsen/logrus"
	"github.com/transip/gotransip/v6"
	"github.com/transip/gotransip/v6/domain"
)

const defaultTTL int = 60 // 60 seconds.

// DNSProvider struct.
type DNSProvider struct {
	configuration *settings.Settings
	clientConfig  gotransip.ClientConfiguration
}

// Init passes DNS settings and store it to the provider instance.
func (provider *DNSProvider) Init(conf *settings.Settings) {
	provider.configuration = conf
	provider.clientConfig = gotransip.ClientConfiguration{
		AccountName: conf.Email}
	if strings.HasPrefix(conf.LoginToken, "-----BEGIN PRIVATE KEY-----") { // Private Key.
		provider.clientConfig.PrivateKeyReader = strings.NewReader(conf.LoginToken)
	} else { // JWT.
		provider.clientConfig.Token = conf.LoginToken
	}
}

// UpdateIP updates the IP address of the given subdomain.
func (provider *DNSProvider) UpdateIP(domainName, subDomainName, ip string) error {
	client, err := gotransip.NewClient(provider.clientConfig)
	if err != nil {
		return err
	}
	domainRepo := domain.Repository{Client: client}

	exists, ttl, err := checkExistence(domainRepo, subDomainName, domainName)
	if err != nil {
		return err
	}

	if exists { // Update.
		err = domainRepo.UpdateDNSEntry(domainName, domain.DNSEntry{
			Name:    subDomainName,
			Type:    provider.setType(),
			Content: ip,
			Expire:  ttl})
		if err != nil {
			log.Error("Failed to update sub domain:", subDomainName)
			return err
		}
	} else { // Create.
		err = domainRepo.AddDNSEntry(domainName, domain.DNSEntry{
			Name:    subDomainName,
			Type:    provider.setType(),
			Content: ip,
			Expire:  defaultTTL})
		if err != nil {
			log.Error("Failed to add sub domain:", subDomainName)
			return err
		}
	}
	return nil
}

func checkExistence(repo domain.Repository, subdomain, domainName string) (bool, int, error) {
	records, err := repo.GetDNSEntries(domainName)
	if err == nil {
		for _, record := range records {
			if record.Name == subdomain {
				return true, record.Expire, nil
			}
		}
	}
	log.Error("Failed to get domain:", domainName)
	return false, defaultTTL, err
}

// defaults to A record (ipv4).
func (provider *DNSProvider) setType() string {
	if strings.ToUpper(provider.configuration.IPType) == utils.IPV6 {
		return utils.IPTypeAAAA
	}
	return utils.IPTypeA
}
