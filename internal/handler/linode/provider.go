package linode

import (
	"context"
	"fmt"

	"github.com/TimothyYe/godns/internal/settings"
	"github.com/TimothyYe/godns/internal/utils"
	"github.com/linode/linodego"
	log "github.com/sirupsen/logrus"
)

type DNSProvider interface {
	Init(conf *settings.Settings) error
	UpdateIP(domainName string, subdomainName string, ip string) error
}

type LinodeDNSProvider struct {
	linodeClient *linodego.Client
}

func (provider *LinodeDNSProvider) Init(conf *settings.Settings) error {
	httpClient, err := CreateHTTPClient(conf)
	if err != nil {
		panic(err)
	}
	linodeAPIClient := linodego.NewClient(httpClient)
	linodeAPIClient.SetDebug(conf.DebugInfo)
	provider.linodeClient = &linodeAPIClient
	return nil
}

func (provider *LinodeDNSProvider) UpdateIP(domain string, subdomain string, ip string) error {
	if subdomain == utils.RootDomain {
		subdomain = ""
	}

	domainID, err := provider.getDomainID(domain)
	if err != nil {
		return err
	}

	recordExists, recordID, err := provider.getDomainRecordID(domainID, subdomain)
	if err != nil {
		return err
	}
	if !recordExists {
		recordID, _ = provider.createDomainRecord(domainID, subdomain)
	}

	err = provider.updateDomainRecord(domainID, recordID, ip)
	if err != nil {
		return err
	}

	return nil
}

func (provider *LinodeDNSProvider) getDomainID(name string) (int, error) {
	f := linodego.Filter{}
	f.AddField(linodego.Eq, "domain", name)
	fStr, err := f.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}

	opts := linodego.NewListOptions(0, string(fStr))
	res, err := provider.linodeClient.ListDomains(context.Background(), opts)
	if err != nil {
		return 0, err
	}
	if len(res) == 0 {
		return 0, fmt.Errorf("No domains found for name %s", name)
	}
	return res[0].ID, nil
}

func (provider *LinodeDNSProvider) getDomainRecordID(domainID int, name string) (bool, int, error) {
	res, err := provider.linodeClient.ListDomainRecords(context.Background(), domainID, nil)
	if err != nil {
		return false, 0, err
	}
	if len(res) == 0 {
		return false, 0, nil
	}
	for _, record := range res {
		if record.Name == name {
			return true, record.ID, nil
		}
	}
	return false, 0, nil
}

func (provider *LinodeDNSProvider) createDomainRecord(domainID int, name string) (int, error) {
	opts := &linodego.DomainRecordCreateOptions{
		Type:   "A",
		Name:   name,
		Target: "127.0.0.1",
		TTLSec: 30,
	}
	record, err := provider.linodeClient.CreateDomainRecord(context.Background(), domainID, *opts)
	if err != nil {
		return 0, err
	}
	return record.ID, nil
}

func (provider *LinodeDNSProvider) updateDomainRecord(domainID int, id int, ip string) error {
	opts := &linodego.DomainRecordUpdateOptions{Target: ip}
	_, err := provider.linodeClient.UpdateDomainRecord(context.Background(), domainID, id, *opts)
	if err != nil {
		return err
	}
	return nil
}
