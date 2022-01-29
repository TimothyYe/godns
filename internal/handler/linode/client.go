package linode

import (
	"context"
	"fmt"

	"github.com/linode/linodego"
	log "github.com/sirupsen/logrus"
)

type LinodeDNSClient struct {
	apiClient *linodego.Client
}

func CreateLinodeDNSClient(linodeClient *linodego.Client) DNSClient {
	dnsClient := LinodeDNSClient{
		apiClient: linodeClient,
	}
	return &dnsClient
}

func (dnsClient *LinodeDNSClient) GetDomainID(name string) (int, error) {
	f := linodego.Filter{}
	f.AddField(linodego.Eq, "domain", name)
	fStr, err := f.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}

	opts := linodego.NewListOptions(0, string(fStr))
	res, err := dnsClient.apiClient.ListDomains(context.Background(), opts)
	if err != nil {
		return 0, err
	}
	if len(res) == 0 {
		return 0, fmt.Errorf("No domains found for name %s", name)
	}
	return res[0].ID, nil
}

func (dnsClient *LinodeDNSClient) GetDomainRecordID(domainID int, name string) (bool, int, error) {
	f := linodego.Filter{}
	f.AddField(linodego.Eq, "name", name)
	fStr, err := f.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}
	opts := linodego.NewListOptions(0, string(fStr))
	res, err := dnsClient.apiClient.ListDomainRecords(context.Background(), domainID, opts)
	if err != nil {
		return false, 0, err
	}
	if len(res) == 0 {
		return false, 0, nil
	}
	return true, res[0].ID, nil
}

func (dnsClient *LinodeDNSClient) CreateDomainRecord(domainID int, name string) (int, error) {
	opts := &linodego.DomainRecordCreateOptions{
		Type:   "A",
		Name:   name,
		Target: "127.0.0.1",
		TTLSec: 30,
	}
	record, err := dnsClient.apiClient.CreateDomainRecord(context.Background(), domainID, *opts)
	if err != nil {
		return 0, err
	}
	return record.ID, nil
}

func (dnsClient *LinodeDNSClient) UpdateDomainRecord(domainId int, id int, ip string) error {
	opts := &linodego.DomainRecordUpdateOptions{Target: ip}
	_, err := dnsClient.apiClient.UpdateDomainRecord(context.Background(), domainId, id, *opts)
	if err != nil {
		return err
	}
	return nil
}
