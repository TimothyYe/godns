package dreamhost

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"runtime/debug"
	"strings"
	"time"

	"github.com/TimothyYe/godns/notify"

	"github.com/TimothyYe/godns"
	"github.com/google/uuid"
)

var (
	// DreamhostURL the API address for dreamhost.com
	DreamhostURL = "https://api.dreamhost.com"
)

// Handler struct
type Handler struct {
	Configuration *godns.Settings
}

// SetConfiguration pass dns settings and store it to handler instance
func (handler *Handler) SetConfiguration(conf *godns.Settings) {
	handler.Configuration = conf
}

// DomainLoop the main logic loop
func (handler *Handler) DomainLoop(domain *godns.Domain, panicChan chan<- godns.Domain) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Recovered in %v: %v\n", err, string(debug.Stack()))
			panicChan <- *domain
		}
	}()

	looping := false
	for {
		if looping {
			// Sleep with interval
			log.Printf("Going to sleep, will start next checking in %d seconds...\r\n", handler.Configuration.Interval)
			time.Sleep(time.Second * time.Duration(handler.Configuration.Interval))
		}
		looping = true

		currentIP, err := godns.GetCurrentIP(handler.Configuration)

		if err != nil {
			log.Println("get_currentIP:", err)
			continue
		}
		log.Println("currentIP is:", currentIP)

		for _, subDomain := range domain.SubDomains {
			hostname := subDomain + "." + domain.DomainName
			lastIP, err := godns.ResolveDNS(hostname, handler.Configuration.Resolver, handler.Configuration.IPType)
			if err != nil {
				log.Println(err)
				continue
			}

			//check against currently known IP, if no change, skip update
			if currentIP == lastIP {
				log.Printf("IP is the same as cached one. Skip update.\n")
			} else {
				log.Printf("%s.%s Start to update record IP...\n", subDomain, domain.DomainName)
				handler.UpdateIP(hostname, currentIP, lastIP)

				// Send notification
				notify.GetNotifyManager(handler.Configuration).Send(fmt.Sprintf("%s.%s", subDomain, domain.DomainName), currentIP)
			}
		}
	}

}

// UpdateIP update subdomain with current IP
func (handler *Handler) UpdateIP(hostname, currentIP, lastIP string) {

	handler.updateDNS(lastIP, currentIP, hostname, "remove")
	handler.updateDNS(lastIP, currentIP, hostname, "add")

}

// updateDNS can add or remove DNS records.
func (handler *Handler) updateDNS(dns, ip, hostname, action string) {
	var ipType string
	if handler.Configuration.IPType == "" || strings.ToUpper(handler.Configuration.IPType) == godns.IPV4 {
		ipType = godns.IPTypeA
	} else if strings.ToUpper(handler.Configuration.IPType) == godns.IPV6 {
		ipType = godns.IPTypeAAAA
	}

	// Generates UUID
	uid, _ := uuid.NewRandom()
	values := url.Values{}
	values.Add("record", hostname)
	values.Add("key", handler.Configuration.LoginToken)
	values.Add("type", ipType)
	values.Add("unique_id", uid.String())
	switch action {
	case "remove":
		// Build URL query (remove)
		values.Add("cmd", "dns-remove_record")
		values.Add("value", dns)
	case "add":
		// Build URL query (add)
		values.Add("cmd", "dns-add_record")
		values.Add("value", ip)
	default:
		log.Fatalf("Unknown action %s\n", action)
	}

	client := godns.GetHttpClient(handler.Configuration, handler.Configuration.UseProxy)
	req, _ := http.NewRequest("POST", DreamhostURL, strings.NewReader(values.Encode()))
	req.SetBasicAuth(handler.Configuration.Email, handler.Configuration.Password)

	if handler.Configuration.UserAgent != "" {
		req.Header.Add("User-Agent", handler.Configuration.UserAgent)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Request error...")
		log.Println("Err:", err.Error())
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode == http.StatusOK {
			log.Println("Update IP success:", string(body))
		} else {
			log.Println("Update IP failed:", string(body))
		}
	}
}
