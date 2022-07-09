package cloudflare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/TimothyYe/godns/internal/settings"
	"github.com/TimothyYe/godns/internal/utils"
	log "github.com/sirupsen/logrus"
)

const (
	// URL is the endpoint for the Cloudflare API.
	URL = "https://api.cloudflare.com/client/v4"
)

// DNSProvider struct definition.
type DNSProvider struct {
	configuration *settings.Settings
	API           string
}

// DNSRecordResponse struct.
type DNSRecordResponse struct {
	Records []DNSRecord `json:"result"`
	Success bool        `json:"success"`
}

// DNSRecordUpdateResponse struct.
type DNSRecordUpdateResponse struct {
	Record  DNSRecord `json:"result"`
	Success bool      `json:"success"`
}

// DNSRecord for Cloudflare API.
type DNSRecord struct {
	ID      string `json:"id"`
	IP      string `json:"content"`
	Name    string `json:"name"`
	Proxied bool   `json:"proxied"`
	Type    string `json:"type"`
	ZoneID  string `json:"zone_id"`
	TTL     int32  `json:"ttl"`
}

// SetIP updates DNSRecord.IP.
func (r *DNSRecord) SetIP(ip string) {
	r.IP = ip
}

// ZoneResponse is a wrapper for Zones.
type ZoneResponse struct {
	Zones   []Zone `json:"result"`
	Success bool   `json:"success"`
}

// Zone object with id and name.
type Zone struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Init passes DNS settings and store it to the provider instance.
func (provider *DNSProvider) Init(conf *settings.Settings) {
	provider.configuration = conf
	provider.API = URL
}

func (provider *DNSProvider) UpdateIP(domainName, subdomainName, ip string) error {
	log.Infof("Checking IP for domain %s", domainName)
	zoneID := provider.getZone(domainName)
	if zoneID != "" {
		records := provider.getDNSRecords(zoneID)

		// update records
		for _, rec := range records {
			if !recordTracked(provider.getCurrentDomain(domainName), &rec) {
				log.Debug("Skipping record:", rec.Name)
				continue
			}
			if rec.IP != ip {
				log.Infof("IP mismatch: Current(%+v) vs Cloudflare(%+v)", ip, rec.IP)
				provider.updateRecord(rec, ip)
			} else {
				log.Infof("Record OK: %+v - %+v", rec.Name, rec.IP)
			}
		}
	} else {
		log.Errorf("Failed to find zone for domain: %s", domainName)
		return fmt.Errorf("failed to find zone for domain: %s", domainName)
	}

	return nil
}

func (provider *DNSProvider) getCurrentDomain(domainName string) *settings.Domain {
	for _, domain := range provider.configuration.Domains {
		if domain.DomainName == domainName {
			return &domain
		}
	}

	return nil
}

// Check if record is present in domain conf.
func recordTracked(domain *settings.Domain, record *DNSRecord) bool {
	for _, subDomain := range domain.SubDomains {
		sd := fmt.Sprintf("%s.%s", subDomain, domain.DomainName)
		if record.Name == sd {
			return true
		} else if subDomain == utils.RootDomain && record.Name == domain.DomainName {
			return true
		}
	}

	return false
}

// Create a new request with auth in place and optional proxy.
func (provider *DNSProvider) newRequest(method, url string, body io.Reader) (*http.Request, *http.Client) {
	client := utils.GetHTTPClient(provider.configuration)
	if client == nil {
		log.Info("cannot create HTTP client")
	}

	req, _ := http.NewRequest(method, provider.API+url, body)
	req.Header.Set("Content-Type", "application/json")

	if provider.configuration.Email != "" && provider.configuration.Password != "" {
		req.Header.Set("X-Auth-Email", provider.configuration.Email)
		req.Header.Set("X-Auth-Key", provider.configuration.Password)
	} else if provider.configuration.LoginToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", provider.configuration.LoginToken))
	}

	return req, client
}

// Find the correct zone via domain name.
func (provider *DNSProvider) getZone(domain string) string {
	var z ZoneResponse

	req, client := provider.newRequest("GET", fmt.Sprintf("/zones?name=%s", domain), nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Request error:", err)
		return ""
	}

	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &z)
	if err != nil {
		log.Errorf("Decoder error: %+v", err)
		log.Debugf("Response body: %+v", string(body))
		return ""
	}
	if !z.Success {
		log.Infof("Response failed: %+v", string(body))
		return ""
	}

	for _, zone := range z.Zones {
		if zone.Name == domain {
			return zone.ID
		}
	}
	return ""
}

// Get all DNS A records for a zone.
func (provider *DNSProvider) getDNSRecords(zoneID string) []DNSRecord {

	var empty []DNSRecord
	var r DNSRecordResponse
	var recordType string

	if provider.configuration.IPType == "" || strings.ToUpper(provider.configuration.IPType) == utils.IPV4 {
		recordType = utils.IPTypeA
	} else if strings.ToUpper(provider.configuration.IPType) == utils.IPV6 {
		recordType = utils.IPTypeAAAA
	}

	log.Infof("Querying records with type: %s", recordType)
	req, client := provider.newRequest("GET", fmt.Sprintf("/zones/"+zoneID+"/dns_records?type=%s&page=1&per_page=500", recordType), nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Request error:", err)
		return empty
	}

	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &r)
	if err != nil {
		log.Infof("Decoder error: %+v", err)
		log.Debugf("Response body: %+v", string(body))
		return empty
	}
	if !r.Success {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Infof("Response failed: %+v", string(body))
		return empty

	}
	return r.Records
}

// Update DNS A Record with new IP.
func (provider *DNSProvider) updateRecord(record DNSRecord, newIP string) string {

	var r DNSRecordUpdateResponse
	record.SetIP(newIP)
	var lastIP string

	j, _ := json.Marshal(record)
	req, client := provider.newRequest("PUT",
		"/zones/"+record.ZoneID+"/dns_records/"+record.ID,
		bytes.NewBuffer(j),
	)
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Request error:", err)
		return ""
	}

	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &r)
	if err != nil {
		log.Errorf("Decoder error: %+v", err)
		log.Debugf("Response body: %+v", string(body))
		return ""
	}
	if !r.Success {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Infof("Response failed: %+v", string(body))
	} else {
		log.Infof("Record updated: %+v - %+v", record.Name, record.IP)
		lastIP = record.IP
	}
	return lastIP
}
