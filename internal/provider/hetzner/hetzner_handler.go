package hetzner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"net/http"

	"github.com/TimothyYe/godns/internal/settings"
	"github.com/TimothyYe/godns/internal/utils"
	log "github.com/sirupsen/logrus"
)

const (
	// URL the API address for Hetzner.
	BaseURL = "https://dns.hetzner.com/api/v1/"
)

type Record struct {
	Type   string `json:"type"`
	ID     string `json:"id"`
	Name   string `json:"name"`
	Value  string `json:"value"`
	TTL    int64  `json:"ttl"`
	ZoneID string `json:"zone_id"`
}

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
		log.Error("Failed to get ZoneID")
		return err
	}

	record, err := provider.getRecord(subdomainName, zoneID, provider.configuration.IPType)

	if err != nil {
		log.Error("Failed to get Record")
		return err
	}
	record.Value = ip

	err = provider.updateRecord(record)
	if err != nil {
		log.Error("Update of Record failed")
	}
	return err
}
func (provider *DNSProvider) getData(endpoint string, param string, value string) ([]byte, error) {

	req, _ := http.NewRequest("GET", BaseURL+endpoint, nil)

	q := req.URL.Query()
	q.Add(param, value)
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Auth-API-Token", provider.configuration.LoginToken)

	resp, err := provider.client.Do(req)

	if err != nil {
		log.Error("Error in fetching: ", err)
		return nil, err
	}

	respBody, _ := io.ReadAll(resp.Body)
	return respBody, nil
}
func (provider *DNSProvider) putData(endpoint string, location string, body []byte) error {

	req, _ := http.NewRequest("PUT", BaseURL+endpoint+"/"+location, bytes.NewBuffer(body))

	req.Header.Add("Auth-API-Token", provider.configuration.LoginToken)
	req.Header.Add("Content-Type", "application/json")

	resp, err := provider.client.Do(req)

	if err != nil {
		log.Error("Fetch failed")
		return err
	}
	if resp.StatusCode != http.StatusOK {
		log.Error("Got non 200 status code: ", resp.Status)
		return fmt.Errorf("got non 200 status code %s", resp.Status)
	}
	return nil
}
func (provider *DNSProvider) getZoneID(zoneName string) (string, error) {

	type Zone struct {
		ID string `json:"id"`
	}

	type GetAllZonesResponse struct {
		Zones []Zone `json:"zones"`
	}

	respBody, err := provider.getData("zones", "name", zoneName)
	if err != nil {
		return "", err
	}
	response := GetAllZonesResponse{}
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return "", err
	}
	if len(response.Zones) == 0 {
		return "", err
	}
	if len(response.Zones) > 1 {
		return "", err
	}
	return response.Zones[0].ID, nil
}

func (provider *DNSProvider) getRecord(recordName string, zoneID string, Type string) (Record, error) {

	type GetRecordsResult struct {
		Records []Record `json:"records"`
	}
	response := GetRecordsResult{}
	respBody, err := provider.getData("records", "zone_id", zoneID)
	if err != nil {
		return Record{}, err
	}
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return Record{}, err
	}
	if len(response.Records) == 0 {
		log.Error("Zone doesn't have any records")
		return Record{}, fmt.Errorf("zone doesn't have an records")
	}
	outRecord := Record{}
	if Type == "IPv6" {
		Type = utils.IPTypeAAAA
	} else {
		Type = utils.IPTypeA
	}
	found := false

	for _, record := range response.Records {

		if record.Name == recordName && record.Type == Type {
			found = true
			outRecord = record
			break
		}
	}
	if found {
		return outRecord, nil
	}
	return outRecord, fmt.Errorf("no record matching value and type found")

}
func (provider *DNSProvider) updateRecord(record Record) error {
	recordJSON, _ := json.Marshal(record)
	err := provider.putData("records", record.ID, recordJSON)
	return err

}
