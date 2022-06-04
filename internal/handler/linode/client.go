package linode

import (
	"context"
	"fmt"

	"github.com/TimothyYe/godns/internal/utils"
	"github.com/linode/linodego"
	log "github.com/sirupsen/logrus"
)

type DNSClient struct {
	linodeClient *linodego.Client
}

func CreateLinodeDNSClient(linodeClient *linodego.Client) IDNSClient {
	dnsClient := DNSClient{
		linodeClient: linodeClient,
	}
	return &dnsClient
}

func (dnsClient *DNSClient) UpdateDNSRecordIP(domain string, subdomain string, ip string) error {
	if subdomain == utils.RootDomain {
		subdomain = ""
	}

	domainID, err := dnsClient.getDomainID(domain)
	if err != nil {
		return err
	}

	recordExists, recordID, err := dnsClient.getDomainRecordID(domainID, subdomain)
	if err != nil {
		return err
	}
	if !recordExists {
		recordID, _ = dnsClient.createDomainRecord(domainID, subdomain)
	}

	err = dnsClient.updateDomainRecord(domainID, recordID, ip)
	if err != nil {
		return err
	}

	return nil
}

func (dnsClient *DNSClient) getDomainID(name string) (int, error) {
	f := linodego.Filter{}
	f.AddField(linodego.Eq, "domain", name)
	fStr, err := f.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}

	opts := linodego.NewListOptions(0, string(fStr))
	res, err := dnsClient.linodeClient.ListDomains(context.Background(), opts)
	if err != nil {
		return 0, err
	}
	if len(res) == 0 {
		return 0, fmt.Errorf("No domains found for name %s", name)
	}
	return res[0].ID, nil
}

func (dnsClient *DNSClient) getDomainRecordID(domainID int, name string) (bool, int, error) {
	f := linodego.Filter{}
	f.AddField(linodego.Eq, "name", name)
	fStr, err := f.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}
	opts := linodego.NewListOptions(0, string(fStr))
	res, err := dnsClient.linodeClient.ListDomainRecords(context.Background(), domainID, opts)
	if err != nil {
		return false, 0, err
	}
	if len(res) == 0 {
		return false, 0, nil
	}
	return true, res[0].ID, nil
}

func (dnsClient *DNSClient) createDomainRecord(domainID int, name string) (int, error) {
	opts := &linodego.DomainRecordCreateOptions{
		Type:   "A",
		Name:   name,
		Target: "127.0.0.1",
		TTLSec: 30,
	}
	record, err := dnsClient.linodeClient.CreateDomainRecord(context.Background(), domainID, *opts)
	if err != nil {
		return 0, err
	}
	return record.ID, nil
}

func (dnsClient *DNSClient) updateDomainRecord(domainID int, id int, ip string) error {
	opts := &linodego.DomainRecordUpdateOptions{Target: ip}
	_, err := dnsClient.linodeClient.UpdateDomainRecord(context.Background(), domainID, id, *opts)
	if err != nil {
		return err
	}
	return nil
}
