package porkbun

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
	// URL is the endpoint for the Porkbun API.
	URL = "https://api.porkbun.com/api/json/v3"
)

// DNSProvider struct definition.
type DNSProvider struct {
	configuration *settings.Settings
	API           string
}

// Record represents a DNS record in Porkbun.
type Record struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content string `json:"content"`
	TTL     string `json:"ttl"`
	Prio    string `json:"prio"`
	Notes   string `json:"notes"`
}

// RecordsResponse represents the response from Porkbun API for retrieving records.
type RecordsResponse struct {
	Status  string   `json:"status"`
	Records []Record `json:"records"`
}

// APIResponse represents a generic API response.
type APIResponse struct {
	Status string `json:"status"`
	ID     string `json:"id,omitempty"`
}

// APIRequest represents the base request structure.
type APIRequest struct {
	SecretAPIKey string `json:"secretapikey"`
	APIKey       string `json:"apikey"`
}

// CreateRecordRequest represents a create record request.
type CreateRecordRequest struct {
	APIRequest
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content string `json:"content"`
	TTL     string `json:"ttl,omitempty"`
}

// EditRecordRequest represents an edit record request.
type EditRecordRequest struct {
	APIRequest
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content string `json:"content"`
	TTL     string `json:"ttl,omitempty"`
}

// Init passes DNS settings and store it to the provider instance.
func (provider *DNSProvider) Init(conf *settings.Settings) {
	provider.configuration = conf
	provider.API = URL
}

// UpdateIP updates the DNS record for the given domain and subdomain.
func (provider *DNSProvider) UpdateIP(domainName, subdomainName, ip string) error {
	log.Infof("Checking IP for domain %s.%s", subdomainName, domainName)

	// Get existing records
	records, err := provider.getRecords(domainName)
	if err != nil {
		log.Errorf("Failed to get DNS records for domain %s: %v", domainName, err)
		return err
	}

	// Determine record type
	recordType := utils.IPTypeA
	if provider.configuration.IPType != "" && strings.ToUpper(provider.configuration.IPType) == utils.IPV6 {
		recordType = utils.IPTypeAAAA
	}

	// Find the target record name
	var targetName string
	if subdomainName == utils.RootDomain {
		targetName = domainName
	} else {
		targetName = fmt.Sprintf("%s.%s", subdomainName, domainName)
	}

	// Check if record exists and needs updating
	var existingRecord *Record
	for _, record := range records {
		if record.Name == targetName && record.Type == recordType {
			existingRecord = &record
			break
		}
	}

	if existingRecord != nil {
		// Record exists, check if IP needs updating
		if existingRecord.Content == ip {
			log.Infof("Record OK: %s - %s", existingRecord.Name, existingRecord.Content)
			return nil
		}

		// Update existing record
		log.Infof("IP mismatch: Current(%s) vs Porkbun(%s)", ip, existingRecord.Content)
		if err := provider.editRecord(domainName, existingRecord.ID, targetName, recordType, ip); err != nil {
			log.Errorf("Failed to update DNS record: %v", err)
			return err
		}
		log.Infof("Record updated: %s - %s", targetName, ip)
	} else {
		// Record doesn't exist, create it
		log.Debugf("Record %s not found, will create it.", targetName)
		if err := provider.createRecord(domainName, targetName, recordType, ip); err != nil {
			log.Errorf("Failed to create DNS record: %v", err)
			return err
		}
		log.Infof("Record [%s] created with IP address: %s", targetName, ip)
	}

	return nil
}

// getRecords retrieves all DNS records for a domain.
func (provider *DNSProvider) getRecords(domain string) ([]Record, error) {
	reqBody := APIRequest{
		SecretAPIKey: provider.configuration.Password,
		APIKey:       provider.configuration.LoginToken,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	client := utils.GetHTTPClient(provider.configuration)
	if client == nil {
		return nil, fmt.Errorf("cannot create HTTP client")
	}

	url := fmt.Sprintf("%s/dns/retrieve/%s", provider.API, domain)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var response RecordsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if response.Status != "SUCCESS" {
		return nil, fmt.Errorf("API returned error status: %s", response.Status)
	}

	return response.Records, nil
}

// createRecord creates a new DNS record.
func (provider *DNSProvider) createRecord(domain, name, recordType, content string) error {
	// For Porkbun API, we need just the subdomain part for the name field
	recordName := ""
	if name != domain {
		recordName = strings.TrimSuffix(name, "."+domain)
	}

	reqBody := CreateRecordRequest{
		APIRequest: APIRequest{
			SecretAPIKey: provider.configuration.Password,
			APIKey:       provider.configuration.LoginToken,
		},
		Name:    recordName,
		Type:    recordType,
		Content: content,
		TTL:     "600", // Default TTL
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}

	client := utils.GetHTTPClient(provider.configuration)
	if client == nil {
		return fmt.Errorf("cannot create HTTP client")
	}

	url := fmt.Sprintf("%s/dns/create/%s", provider.API, domain)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	var response APIResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if response.Status != "SUCCESS" {
		return fmt.Errorf("API returned error status: %s", response.Status)
	}

	return nil
}

// editRecord updates an existing DNS record.
func (provider *DNSProvider) editRecord(domain, recordID, name, recordType, content string) error {
	// For Porkbun API, we need just the subdomain part for the name field
	recordName := ""
	if name != domain {
		recordName = strings.TrimSuffix(name, "."+domain)
	}

	reqBody := EditRecordRequest{
		APIRequest: APIRequest{
			SecretAPIKey: provider.configuration.Password,
			APIKey:       provider.configuration.LoginToken,
		},
		Name:    recordName,
		Type:    recordType,
		Content: content,
		TTL:     "600", // Default TTL
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}

	client := utils.GetHTTPClient(provider.configuration)
	if client == nil {
		return fmt.Errorf("cannot create HTTP client")
	}

	url := fmt.Sprintf("%s/dns/edit/%s/%s", provider.API, domain, recordID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	var response APIResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if response.Status != "SUCCESS" {
		return fmt.Errorf("API returned error status: %s", response.Status)
	}

	return nil
}
