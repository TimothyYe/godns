package cloudflare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/TimothyYe/godns/internal/settings"
	"github.com/TimothyYe/godns/internal/utils"
	"github.com/TimothyYe/godns/pkg/notification"

	log "github.com/sirupsen/logrus"
)

// Handler struct definition.
type Handler struct {
	Configuration *settings.Settings
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

// SetConfiguration pass dns settings and store it to handler instance.
func (handler *Handler) SetConfiguration(conf *settings.Settings) {
	handler.Configuration = conf
	handler.API = "https://api.cloudflare.com/client/v4"
}

// DomainLoop the main logic loop.
func (handler *Handler) DomainLoop(domain *settings.Domain, panicChan chan<- settings.Domain, runOnce bool) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("Recovered in %v: %v", err, string(debug.Stack()))
			panicChan <- *domain
		}
	}()

	var lastIP string
	looping := false
	for while := true; while; while = !runOnce {
		if looping {
			// Sleep with interval
			log.Debugf("Going to sleep, will start next checking in %d seconds...", handler.Configuration.Interval)
			time.Sleep(time.Second * time.Duration(handler.Configuration.Interval))
		}
		looping = true

		currentIP, err := utils.GetCurrentIP(handler.Configuration)
		if err != nil {
			log.Error("Error in GetCurrentIP:", err)
			continue
		}
		log.Debug("Current IP is:", currentIP)
		//check against locally cached IP, if no change, skip update
		if currentIP == lastIP {
			log.Infof("IP is the same as cached one (%s). Skip update.", currentIP)
		} else {
			log.Infof("Checking IP for domain %s", domain.DomainName)
			zoneID := handler.getZone(domain.DomainName)
			if zoneID != "" {
				records := handler.getDNSRecords(zoneID)

				// update records
				for _, rec := range records {
					if !recordTracked(domain, &rec) {
						log.Debug("Skiping record:", rec.Name)
						continue
					}
					if rec.IP != currentIP {
						log.Infof("IP mismatch: Current(%+v) vs Cloudflare(%+v)", currentIP, rec.IP)
						lastIP = handler.updateRecord(rec, currentIP)

						// Send notification
						notification.GetNotificationManager(handler.Configuration).Send(rec.Name, currentIP)
					} else {
						log.Infof("Record OK: %+v - %+v", rec.Name, rec.IP)
					}
				}
			} else {
				log.Infof("Failed to find zone for domain: %s", domain.DomainName)
			}
		}
	}
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
func (handler *Handler) newRequest(method, url string, body io.Reader) (*http.Request, *http.Client) {
	client := utils.GetHTTPClient(handler.Configuration, handler.Configuration.UseProxy)
	if client == nil {
		log.Info("cannot create HTTP client")
	}

	req, _ := http.NewRequest(method, handler.API+url, body)
	req.Header.Set("Content-Type", "application/json")

	if handler.Configuration.Email != "" && handler.Configuration.Password != "" {
		req.Header.Set("X-Auth-Email", handler.Configuration.Email)
		req.Header.Set("X-Auth-Key", handler.Configuration.Password)
	} else if handler.Configuration.LoginToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", handler.Configuration.LoginToken))
	}

	return req, client
}

// Find the correct zone via domain name.
func (handler *Handler) getZone(domain string) string {

	var z ZoneResponse

	req, client := handler.newRequest("GET", fmt.Sprintf("/zones?name=%s", domain), nil)
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
func (handler *Handler) getDNSRecords(zoneID string) []DNSRecord {

	var empty []DNSRecord
	var r DNSRecordResponse
	var recordType string

	if handler.Configuration.IPType == "" || strings.ToUpper(handler.Configuration.IPType) == utils.IPV4 {
		recordType = utils.IPTypeA
	} else if strings.ToUpper(handler.Configuration.IPType) == utils.IPV6 {
		recordType = utils.IPTypeAAAA
	}

	log.Infof("Querying records with type: %s", recordType)
	req, client := handler.newRequest("GET", fmt.Sprintf("/zones/"+zoneID+"/dns_records?type=%s&page=1&per_page=500", recordType), nil)
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
func (handler *Handler) updateRecord(record DNSRecord, newIP string) string {

	var r DNSRecordUpdateResponse
	record.SetIP(newIP)
	var lastIP string

	j, _ := json.Marshal(record)
	req, client := handler.newRequest("PUT",
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
