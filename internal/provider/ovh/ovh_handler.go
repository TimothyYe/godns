package ovh

import (
	"fmt"
	"strings"

	"github.com/TimothyYe/godns/internal/settings"
	"github.com/TimothyYe/godns/internal/utils"
	"github.com/ovh/go-ovh/ovh"
	log "github.com/sirupsen/logrus"
)

type DNSProvider struct {
	configuration *settings.Settings
}

func (provider *DNSProvider) Init(conf *settings.Settings) {
	provider.configuration = conf
}

type Record struct {
	Zone      string `json:"zone"`
	TTL       int    `json:"ttl"`
	Value     string `json:"target"`
	SubDomain string `json:"subDomain"`
	Type      string `json:"fieldType"`
	ID        int    `json:"id"`
}

func (provider *DNSProvider) UpdateIP(domainName string, subdomainName string, ip string) error {
	client, err := ovh.NewClient(
		"ovh-eu",
		provider.configuration.AppKey,
		provider.configuration.AppSecret,
		provider.configuration.ConsumerKey,
	)
	if err != nil {
		log.Error("OVH Client error: ", err)
		return err
	}
	var IDs []int
	query := fmt.Sprintf("/domain/zone/%s/record?subDomain=%s", domainName, subdomainName)

	err = client.Get(query, &IDs)
	if err != nil {
		log.Error("Fetch error")
		return err
	}
	if len(IDs) < 1 {
		log.Error("No matching records")
		return fmt.Errorf("no matching records")
	}
	outrec := Record{}

	for _, id := range IDs {
		record := Record{}
		err = client.Get(fmt.Sprintf("/domain/zone/%s/record/%s", domainName, fmt.Sprint(id)), &record)
		if err != nil {
			log.Error("Fetch error on get record: ", id)
			return err
		}

		if strings.ToUpper(provider.configuration.IPType) == provider.recordTypeToIPType(record.Type) {
			outrec = record
			break
		}
	}
	if outrec.ID == 0 {
		log.Error("No fitting record type found")
		return fmt.Errorf("no fitting record type found")
	}
	outrec.Value = ip
	// Update IP.
	err = client.Put(fmt.Sprintf("/domain/zone/%s/record/%s", domainName, fmt.Sprint(outrec.ID)), outrec, nil)
	if err != nil {
		log.Error("Error while Updating record: ", outrec.SubDomain, outrec.Zone)
		return err
	}
	// Refresh zone.
	err = client.Post(fmt.Sprintf("/domain/zone/%s/refresh", domainName), nil, nil)
	if err != nil {
		log.Error("Applying new records failed")
		return err
	}
	return nil
}
func (provider *DNSProvider) recordTypeToIPType(Type string) string {
	if Type == utils.IPTypeAAAA {
		return utils.IPV6
	}
	return utils.IPV4

}
