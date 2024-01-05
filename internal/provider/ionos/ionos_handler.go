package ionos

// API Docs: https://developer.hosting.ionos.com/docs/dns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/TimothyYe/godns/internal/settings"
	"github.com/TimothyYe/godns/internal/utils"
	"github.com/sirupsen/logrus"
)

const (
	BaseURL = "https://api.hosting.ionos.com/dns/v1/"
)

// DNSProvider struct.
type DNSProvider struct {
	configuration *settings.Settings
	client        *http.Client
}

// Init passes DNS settings and store it to the provider instance.
func (provider *DNSProvider) Init(conf *settings.Settings) {
	provider.configuration = conf
	provider.client = utils.GetHTTPClient(provider.configuration)
}

func (provider *DNSProvider) UpdateIP(domainName, subdomainName, ip string) error {
	zoneID, err := provider.getZoneID(domainName)
	if err != nil {
		return err
	}

	recordID, currIP, err := provider.getRecord(zoneID, subdomainName+"."+domainName)
	if err != nil {
		return err
	}

	if currIP == ip {
		return nil
	}

	return provider.updateRecord(zoneID, recordID, subdomainName+"."+domainName, ip)
}

func (provider *DNSProvider) getData(endpoint string, params map[string]string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, BaseURL+endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-API-Key", provider.configuration.LoginToken)

	if params != nil {
		q := req.URL.Query()
		for k, v := range params {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	resp, err := provider.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get data from %s, status code: %s", BaseURL+endpoint, resp.Status)
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)

}

func (provider *DNSProvider) putData(endpoint string, params map[string]any) error {

	var body []byte
	var err error
	if params != nil {
		body, err = json.Marshal(params)
		if err != nil {
			return err
		}
	}

	req, err := http.NewRequest(http.MethodPut, BaseURL+endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-API-Key", provider.configuration.LoginToken)

	resp, err := provider.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to PUT %s, status: %s", endpoint, resp.Status)
	}
	defer resp.Body.Close()

	return nil
}

type zoneResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func (provider *DNSProvider) getZoneID(domainName string) (string, error) {

	body, err := provider.getData("zones", nil)
	if err != nil {
		return "", err
	}

	var zones []zoneResponse
	err = json.Unmarshal(body, &zones)
	if err != nil {
		return "", err
	}

	for _, zone := range zones {
		if zone.Name == domainName {
			return zone.ID, nil
		}
	}

	return "", fmt.Errorf("zone %s not found", domainName)
}

type recordResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	RootName string `json:"rootName"`
	Type     string `json:"type"`
	Content  string `json:"content"`
	TTL      int    `json:"ttl"`
	Prio     int    `json:"prio"`
	Disabled bool   `json:"disabled"`
}

type recordListResponse struct {
	zoneResponse
	Records []recordResponse `json:"records"`
}

func (provider *DNSProvider) getRecord(zoneID, recordName string) (id string, ip string, err error) {

	ipType := utils.IPTypeA
	if provider.configuration.IPType == utils.IPV6 || provider.configuration.IPType == utils.IPTypeAAAA {
		ipType = utils.IPTypeAAAA
	}

	body, err := provider.getData(fmt.Sprintf("zones/%s", zoneID),
		map[string]string{
			"recordName": recordName,
			"recordType": ipType,
		})
	if err != nil {
		return "", "", err
	}

	var rlp recordListResponse
	err = json.Unmarshal(body, &rlp)
	if err != nil {
		return "", "", err
	}

	if len(rlp.Records) > 0 {
		return rlp.Records[0].ID, rlp.Records[0].Content, nil
	}

	return "", "", fmt.Errorf("record %s not found", recordName)
}

func (provider *DNSProvider) updateRecord(zoneID, recordID, recordName, ip string) error {

	err := provider.putData(fmt.Sprintf("zones/%s/records/%s", zoneID, recordID), map[string]any{"content": ip})
	if err != nil {
		return fmt.Errorf("failed to update record %s: %w", recordName, err)
	}

	logrus.Infof("Updated record %s to %s", recordName, ip)

	return nil
}
