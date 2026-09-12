package namecom

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
	// BaseURL is the endpoint of the name.com CORE v1 API.
	BaseURL = "https://api.name.com/core/v1"
	// DefaultTTL is used for newly created records. name.com allows a minimum TTL of 300 seconds.
	DefaultTTL = 300
	// pageSize is the number of records requested per page when listing a zone.
	pageSize = 500
)

// DNSProvider struct definition.
type DNSProvider struct {
	configuration *settings.Settings
	client        *http.Client
	// API is the base address of the name.com API. It is overridden in tests.
	API string
}

// Record is an individual DNS resource record as returned by the name.com API.
type Record struct {
	ID         int64  `json:"id,omitempty"`
	DomainName string `json:"domainName,omitempty"`
	Host       string `json:"host"`
	FQDN       string `json:"fqdn,omitempty"`
	Type       string `json:"type"`
	Answer     string `json:"answer"`
	TTL        int64  `json:"ttl"`
	Priority   int64  `json:"priority,omitempty"`
}

// ListRecordsResponse is the response of the ListRecords endpoint.
type ListRecordsResponse struct {
	Records  []Record `json:"records"`
	NextPage int      `json:"nextPage"`
}

// recordBody is the payload sent to CreateRecord and UpdateRecord.
type recordBody struct {
	Host   string `json:"host"`
	Type   string `json:"type"`
	Answer string `json:"answer"`
	TTL    int64  `json:"ttl"`
}

// errorResponse is the body returned by the API on failure.
type errorResponse struct {
	Message string `json:"message"`
	Details string `json:"details"`
}

// Init passes DNS settings and store it to the provider instance.
func (provider *DNSProvider) Init(conf *settings.Settings) {
	provider.configuration = conf
	provider.client = utils.GetHTTPClient(conf)
	provider.API = BaseURL
}

// UpdateIP updates the DNS record for the given domain and subdomain.
func (provider *DNSProvider) UpdateIP(domainName, subdomainName, ip string) error {
	log.Infof("Checking IP for domain %s.%s", subdomainName, domainName)

	records, err := provider.listRecords(domainName)
	if err != nil {
		log.Errorf("Failed to get DNS records for domain %s: %v", domainName, err)
		return err
	}

	recordType := utils.RecordType(provider.configuration.IPType)
	host := hostName(subdomainName)

	var existing *Record
	for i := range records {
		if hostName(records[i].Host) == host && records[i].Type == recordType {
			existing = &records[i]
			break
		}
	}

	if existing == nil {
		log.Debugf("Record %s.%s not found, will create it.", subdomainName, domainName)
		if err := provider.createRecord(domainName, host, recordType, ip); err != nil {
			log.Errorf("Failed to create DNS record: %v", err)
			return err
		}
		log.Infof("Record [%s.%s] created with IP address: %s", subdomainName, domainName, ip)
		return nil
	}

	if existing.Answer == ip {
		log.Infof("Record OK: %s.%s - %s", subdomainName, domainName, ip)
		return nil
	}

	log.Infof("IP mismatch: Current(%s) vs name.com(%s)", ip, existing.Answer)
	if err := provider.updateRecord(domainName, existing, ip); err != nil {
		log.Errorf("Failed to update DNS record: %v", err)
		return err
	}
	log.Infof("Record updated: %s.%s - %s", subdomainName, domainName, ip)

	return nil
}

// hostName normalizes a host so that the apex can be compared regardless of
// whether it is spelled "@" or "" (name.com accepts and returns both).
func hostName(host string) string {
	if host == "" {
		return utils.RootDomain
	}
	return host
}

// listRecords retrieves every DNS record of a zone, following pagination.
func (provider *DNSProvider) listRecords(domainName string) ([]Record, error) {
	var records []Record
	page := 1

	for {
		url := fmt.Sprintf("%s/domains/%s/records?perPage=%d&page=%d", provider.API, domainName, pageSize, page)
		body, err := provider.doRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}

		var response ListRecordsResponse
		if err := json.Unmarshal(body, &response); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}
		records = append(records, response.Records...)

		if response.NextPage == 0 || response.NextPage <= page {
			return records, nil
		}
		page = response.NextPage
	}
}

// createRecord creates a new DNS record in the zone.
func (provider *DNSProvider) createRecord(domainName, host, recordType, ip string) error {
	url := fmt.Sprintf("%s/domains/%s/records", provider.API, domainName)
	payload := recordBody{
		Host:   host,
		Type:   recordType,
		Answer: ip,
		TTL:    DefaultTTL,
	}

	_, err := provider.doRequest(http.MethodPost, url, payload)
	return err
}

// updateRecord replaces an existing DNS record. The name.com API treats PUT as a
// full overwrite, so every required field is resent with the new answer.
func (provider *DNSProvider) updateRecord(domainName string, record *Record, ip string) error {
	url := fmt.Sprintf("%s/domains/%s/records/%d", provider.API, domainName, record.ID)
	ttl := record.TTL
	if ttl < DefaultTTL {
		ttl = DefaultTTL
	}
	payload := recordBody{
		Host:   record.Host,
		Type:   record.Type,
		Answer: ip,
		TTL:    ttl,
	}

	_, err := provider.doRequest(http.MethodPut, url, payload)
	return err
}

// doRequest sends an authenticated request to the API and returns the response body.
func (provider *DNSProvider) doRequest(method, url string, payload interface{}) ([]byte, error) {
	var reqBody io.Reader
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request: %w", err)
		}
		reqBody = bytes.NewBuffer(data)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.SetBasicAuth(provider.configuration.Email, provider.configuration.LoginToken)
	req.Header.Set("Accept", "application/json")
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := provider.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErr errorResponse
		if json.Unmarshal(body, &apiErr) == nil && apiErr.Message != "" {
			if apiErr.Details != "" {
				return nil, fmt.Errorf("name.com API error (%s): %s: %s", resp.Status, apiErr.Message, apiErr.Details)
			}
			return nil, fmt.Errorf("name.com API error (%s): %s", resp.Status, apiErr.Message)
		}
		return nil, fmt.Errorf("name.com API error (%s): %s", resp.Status, string(body))
	}

	return body, nil
}
