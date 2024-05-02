package digitalocean

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/TimothyYe/godns/internal/settings"
	"github.com/TimothyYe/godns/internal/utils"
	log "github.com/sirupsen/logrus"
)

const (
	// URL is the endpoint for the DigitalOcean API.
	URL = "https://api.digitalocean.com/v2"
)

// DNSProvider struct definition.
type DNSProvider struct {
	configuration *settings.Settings
	API           string
}

type DomainRecordsResponse struct {
	Records []DNSRecord `json:"domain_records"`
}

// DNSRecord for DigitalOcean API.
type DNSRecord struct {
	ID   int32  `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
	IP   string `json:"data"`
	TTL  int32  `json:"ttl"`
}

// SetIP updates DNSRecord.IP.
func (r *DNSRecord) SetIP(ip string) {
	r.IP = ip
}

// Init passes DNS settings and store it to the provider instance.
func (provider *DNSProvider) Init(conf *settings.Settings) {
	provider.configuration = conf
	provider.API = URL
}

func (provider *DNSProvider) UpdateIP(domainName, subdomainName, ip string) error {
	log.Infof("Checking IP for domain %s", domainName)

	records := provider.getDNSRecords(domainName)
	matched := false

	// update records
	for _, rec := range records {
		rec := rec
		if !recordTracked(provider.getCurrentDomain(domainName), &rec) {
			log.Debug("Skipping record:", rec.Name)
			continue
		}

		if strings.Contains(rec.Name, subdomainName) || rec.Name == domainName {
			if rec.IP != ip {
				log.Infof("IP mismatch: Current(%+v) vs DigitalOcean(%+v)", ip, rec.IP)
				provider.updateRecord(domainName, rec, ip)
			} else {
				log.Infof("Record OK: %+v - %+v", rec.Name, rec.IP)
			}

			matched = true
		}
	}

	if !matched {
		log.Debugf("Record %s not found, will create it.", subdomainName)
		if err := provider.createRecord(domainName, subdomainName, ip); err != nil {
			return err
		}
		log.Infof("Record [%s] created with IP address: %s", subdomainName, ip)
	}

	return nil
}

func (provider *DNSProvider) getRecordType() string {
	var recordType string = utils.IPTypeA
	if provider.configuration.IPType == "" || strings.ToUpper(provider.configuration.IPType) == utils.IPV4 {
		recordType = utils.IPTypeA
	} else if strings.ToUpper(provider.configuration.IPType) == utils.IPV6 {
		recordType = utils.IPTypeAAAA
	}

	return recordType
}

func (provider *DNSProvider) getCurrentDomain(domainName string) *settings.Domain {
	for _, domain := range provider.configuration.Domains {
		domain := domain
		if domain.DomainName == domainName {
			return &domain
		}
	}

	return nil
}

// Check if record is present in domain conf.
func recordTracked(domain *settings.Domain, record *DNSRecord) bool {
	for _, subDomain := range domain.SubDomains {
		if record.Name == subDomain {
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
	log.Debugf("Created %+v request for %+v", string(method), string(url))

	return req, client
}

// Get all DNS A(AAA) records for a zone.
func (provider *DNSProvider) getDNSRecords(domainName string) []DNSRecord {

	var empty []DNSRecord
	var r DomainRecordsResponse
	recordType := provider.getRecordType()

	log.Infof("Querying records with type: %s", recordType)
	req, client := provider.newRequest("GET", fmt.Sprintf("/domains/"+domainName+"/records?type=%s&page=1&per_page=200", recordType), nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Request error:", err)
		return empty
	}

	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &r)
	if err != nil {
		log.Infof("Decoder error: %+v", err)
		log.Debugf("Response body: %+v", string(body))
		return empty
	}

	return r.Records
}

func (provider *DNSProvider) createRecord(domain, subDomain, ip string) error {
	recordType := provider.getRecordType()

	newRecord := DNSRecord{
		Type: recordType,
		IP:   ip,
		TTL:  int32(provider.configuration.Interval),
	}

	if subDomain == utils.RootDomain {
		newRecord.Name = utils.RootDomain
	} else {
		newRecord.Name = subDomain
	}

	content, err := json.Marshal(newRecord)
	if err != nil {
		log.Errorf("Encoder error: %+v", err)
		return err
	}

	req, client := provider.newRequest("POST", fmt.Sprintf("/domains/%s/records", domain), bytes.NewBuffer(content))
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Request error:", err)
		return err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Failed to read request body: %+v", err)
		return err
	}

	var r DNSRecord
	err = json.Unmarshal(body, &r)
	if err != nil {
		log.Errorf("Response decoder error: %+v", err)
		log.Debugf("Response body: %+v", string(body))
		return err
	}

	return nil
}

// Update DNS Record with new IP.
func (provider *DNSProvider) updateRecord(domainName string, record DNSRecord, newIP string) string {

	var r DNSRecord
	record.SetIP(newIP)
	var lastIP string

	j, _ := json.Marshal(record)
	req, client := provider.newRequest("PUT",
		fmt.Sprintf("/domains/%s/records/%d", domainName, record.ID),
		bytes.NewBuffer(j),
	)
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Request error:", err)
		return ""
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &r)
	if err != nil {
		log.Errorf("Decoder error: %+v", err)
		log.Debugf("Response body: %+v", string(body))
		return ""
	}
	log.Infof("Record updated: %+v - %+v", record.Name, record.IP)
	lastIP = record.IP

	return lastIP
}
