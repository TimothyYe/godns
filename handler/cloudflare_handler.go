package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/TimothyYe/godns"
	"golang.org/x/net/proxy"
)

// CloudflareHandler struct definition
type CloudflareHandler struct {
	Configuration *godns.Settings
	API           string
}

// DNS api response
type DNSRecordResponse struct {
	Records []DNSRecord `json:"result"`
	Success bool        `json:"success"`
}

// DNS update api response
type DNSRecordUpdateResponse struct {
	Record  DNSRecord `json:"result"`
	Success bool      `json:"success"`
}

type DNSRecord struct {
	Id      string `json:"id"`
	Ip      string `json:"content"`
	Name    string `json:"name"`
	Proxied bool   `json:"proxied"`
	Type    string `json:"type"`
	ZoneId  string `json:"zone_id"`
}

func (r *DNSRecord) SetIp(ip string) {
	r.Ip = ip
}

// response from zone api request
type ZoneResponse struct {
	Zones   []Zone `json:"result"`
	Success bool   `json:"success"`
}

// nested results, only care about name and id
type Zone struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// SetConfiguration pass dns settings and store it to handler instance
func (handler *CloudflareHandler) SetConfiguration(conf *godns.Settings) {
	handler.Configuration = conf
	handler.API = "https://api.cloudflare.com/client/v4"
}

// DomainLoop the main logic loop
func (handler *CloudflareHandler) DomainLoop(domain *godns.Domain, panicChan chan<- godns.Domain) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Recovered in %v: %v\n", err, debug.Stack())
			panicChan <- *domain
		}
	}()

	for {
		currentIp, err := godns.GetCurrentIP(handler.Configuration)
		if err != nil {
			log.Println("Error in GetCurrentIP:", err)
			continue
		}
		log.Println("Current IP is:", currentIp)
		// TODO: check against locally cached IP, if no change, skip update

		log.Println("Checking IP for domain", domain.DomainName)
		zoneId := handler.getZone(domain.DomainName)
		if zoneId != "" {
			records := handler.getDNSRecords(zoneId)

			// update records
			for _, rec := range records {
				if recordTracked(domain, &rec) != true {
					log.Println("Skiping record:", rec.Name)
					continue
				}
				if rec.Ip != currentIp {
					log.Printf("IP mismatch: Current(%+v) vs Cloudflare(%+v)\r\n", currentIp, rec.Ip)
					handler.updateRecord(rec, currentIp)
				} else {
					log.Printf("Record OK: %+v - %+v\r\n", rec.Name, rec.Ip)
				}
			}
		} else {
			log.Println("Failed to find zone for domain:", domain.DomainName)
		}

		// Interval is 5 minutes
		log.Printf("Going to sleep, will start next checking in %d minutes...\r\n", godns.INTERVAL)
		time.Sleep(time.Minute * godns.INTERVAL)
	}
}

// Check if record is present in domain conf
func recordTracked(domain *godns.Domain, record *DNSRecord) bool {

	if record.Name == domain.DomainName {
		return true
	}

	for _, subDomain := range domain.SubDomains {
		sd := subDomain + "." + domain.DomainName
		if record.Name == sd {
			return true
		}
	}

	return false
}

// Create a new request with auth in place and optional proxy
func (handler *CloudflareHandler) newRequest(method, url string, body io.Reader) (*http.Request, *http.Client) {
	client := &http.Client{}

	if handler.Configuration.Socks5Proxy != "" {
		log.Println("use socks5 proxy:" + handler.Configuration.Socks5Proxy)
		dialer, err := proxy.SOCKS5("tcp", handler.Configuration.Socks5Proxy, nil, proxy.Direct)
		if err != nil {
			log.Println("can't connect to the proxy:", err)
		} else {
			httpTransport := &http.Transport{}
			client.Transport = httpTransport
			httpTransport.Dial = dialer.Dial
		}
	}

	req, _ := http.NewRequest(method, handler.API+url, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Email", handler.Configuration.Email)
	req.Header.Set("X-Auth-Key", handler.Configuration.Password)
	return req, client
}

// Find the correct zone via domain name
func (handler *CloudflareHandler) getZone(domain string) string {

	var z ZoneResponse

	req, client := handler.newRequest("GET", "/zones", nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Request error:", err.Error())
		return ""
	}

	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &z)
	if err != nil {
		log.Printf("Decoder error: %+v\n", err)
		log.Printf("Response body: %+v\n", string(body))
		return ""
	}
	if z.Success != true {
		log.Printf("Response failed: %+v\n", string(body))
		return ""
	}

	for _, zone := range z.Zones {
		if zone.Name == domain {
			return zone.Id
		}
	}
	return ""
}

// Get all DNS A records for a zone
func (handler *CloudflareHandler) getDNSRecords(zoneId string) []DNSRecord {

	var empty []DNSRecord
	var r DNSRecordResponse

	req, client := handler.newRequest("GET", "/zones/"+zoneId+"/dns_records?type=A", nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Request error:", err.Error())
		return empty
	}

	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &r)
	if err != nil {
		log.Printf("Decoder error: %+v\n", err)
		log.Printf("Response body: %+v\n", string(body))
		return empty
	}
	if r.Success != true {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("Response failed: %+v\n", string(body))
		return empty

	}
	return r.Records
}

// Update DNS A Record with new IP
func (handler *CloudflareHandler) updateRecord(record DNSRecord, newIp string) {

	var r DNSRecordUpdateResponse
	record.SetIp(newIp)

	j, _ := json.Marshal(record)
	req, client := handler.newRequest("PUT",
		"/zones/"+record.ZoneId+"/dns_records/"+record.Id,
		bytes.NewBuffer(j),
	)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Request error:", err.Error())
		return
	}

	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &r)
	if err != nil {
		log.Printf("Decoder error: %+v\n", err)
		log.Printf("Response body: %+v\n", string(body))
		return
	}
	if r.Success != true {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("Response failed: %+v\n", string(body))
	} else {
		log.Printf("Record updated: %+v - %+v", record.Name, record.Ip)
	}
}
