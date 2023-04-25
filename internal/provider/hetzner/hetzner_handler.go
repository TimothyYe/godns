package hetzner

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/TimothyYe/godns/internal/settings"
)

const (
	// URL the API address for Strato.
	BASE_URL = "https://dns.hetzner.com/api/v1/"
)

type Record struct {
	Type  string `json:"type`
	Id    string `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
	TTL   int64  `json:"ttl`
}

// DNSProvider struct.
type DNSProvider struct {
	configuration *settings.Settings
}

// Init passes DNS settings and store it to the provider instance.
func (provider *DNSProvider) Init(conf *settings.Settings) {
	provider.configuration = conf
}

func (provider *DNSProvider) UpdateIP(domainName, subdomainName, ip string) error {
	return nil
}
func (provider *DNSProvider) getData(endpoint string, param string) ([]byte, error) {
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("GET", BASE_URL+"/zones", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("name", param)
	req.URL.RawQuery = q.Encode()

	// Headers
	req.Header.Add("Auth-API-Token", provider.configuration.LoginToken)

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	// Read Response Body
	respBody, _ := ioutil.ReadAll(resp.Body)
	return respBody, nil
}
func (provider *DNSProvider) getZoneID(zone_name string) (string, error) {

	type Zone struct {
		Id string `json: "id"`
	}

	type GetAllZonesResponse struct {
		Zones []Zone `json: "zones"`
	}

	// Create client

	response := GetAllZonesResponse{}
	respBody, err := provider.getData("zones", zone_name)
	if err != nil {
		return "", err
	}

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
	return response.Zones[0].Id, nil
	// Display Results

}

func (provider *DNSProvider) getRecordID(record_name string) (Record, error) {

	type GetRecordsResult struct {
		Records []Record `json:"records"`
	}
	response := GetRecordsResult{}
	respBody, err := provider.getData("records", record_name)
	if err != nil {
		return Record{}, err
	}

	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return Record{}, err
	}
	if len(response.Records) == 0 {
		return Record{}, err
	}
	outRecord := Record{}
	found := false
	for _, record := range response.Records {
		if record.Name == record_name {
			found = true
			outRecord = record
			break
		}
	}
	if found {
		return outRecord, nil
	} else {
		return outRecord, err
	}

}
