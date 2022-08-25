package scaleway

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/TimothyYe/godns/internal/settings"
	"github.com/TimothyYe/godns/internal/utils"
	log "github.com/sirupsen/logrus"
)

const (
	URL = "https://api.scaleway.com/domain/v2beta1/dns-zones/%s/records"
)

// DNSProvider struct.
type DNSProvider struct {
	configuration *settings.Settings
}

// Record for Scaleway API.
type Record struct {
	Name    string `json:"name"`
	Data    string `json:"data"`
	TTL     int    `json:"ttl"`
	Comment string `json:"comment"`
}

// IDFields to filter DNS records for Scaleway API.
type IDFields struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// SetRecord for Scaleway API.
type SetRecord struct {
	IDFields IDFields `json:"id_fields"`
	Records  []Record `json:"records"`
}

// DNSChange for Scaleway API.
type DNSChange struct {
	Set SetRecord `json:"set"`
}

// DNSUpdateRequest for Scaleway API.
type DNSUpdateRequest struct {
	Changes []DNSChange `json:"changes"`
}

func (provider *DNSProvider) Init(conf *settings.Settings) {
	provider.configuration = conf
}

func (provider *DNSProvider) UpdateIP(domainName string, subdomainName string, ip string) error {
	log.Infof("%s.%s - Start to update record IP...", subdomainName, domainName)
	err := provider.updateIP(domainName, subdomainName, ip)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (provider *DNSProvider) getRecordType() (string, error) {
	if strings.ToUpper(provider.configuration.IPType) == utils.IPV4 {
		return utils.IPTypeA, nil
	} else if strings.ToUpper(provider.configuration.IPType) == utils.IPV6 {
		return utils.IPTypeAAAA, nil
	} else {
		return "", errors.New("must specify \"ip_type\" in config for Scaleway")
	}
}

// updateIP update subdomain with current IP.
func (provider *DNSProvider) updateIP(domain, subDomain, currentIP string) error {
	recordType, err := provider.getRecordType()
	if err != nil {
		return err
	}

	reqBody := DNSUpdateRequest{Changes: []DNSChange{{SetRecord{
		IDFields: IDFields{
			Name: subDomain,
			Type: recordType,
		},
		Records: []Record{
			{
				Name:    subDomain,
				Data:    currentIP,
				TTL:     provider.configuration.Interval,
				Comment: "Set by GoDNS",
			},
		},
	}}}}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return errors.New("failed to encode request body as json")
	}

	req, _ := http.NewRequest("PATCH", fmt.Sprintf(URL, domain), bytes.NewReader(jsonBody))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Auth-Token", provider.configuration.LoginToken)
	if provider.configuration.UserAgent != "" {
		req.Header.Add("User-Agent", provider.configuration.UserAgent)
	}

	client := utils.GetHTTPClient(provider.configuration)
	log.Debugf("Requesting update for '%s.%s': '%v'", subDomain, domain, reqBody)
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Request error:", err)
		return errors.New("failed to complete update request")
	}

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Debugf("Update failed for '%s.%s': %s", subDomain, domain, string(body))
		return fmt.Errorf("update IP failed with status '%d'", resp.StatusCode)
	}
	log.Debugf("Update IP success for '%s.%s': '%s'", subDomain, domain, string(body))
	return nil
}
