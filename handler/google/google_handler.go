package google

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/TimothyYe/godns"
)

var (
	// GoogleURL the API address for Google Domains
	GoogleURL = "https://%s:%s@domains.google.com/nic/update?hostname=%s.%s&myip=%s"
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
			log.Printf("Recovered in %v: %v\n", err, debug.Stack())
			panicChan <- *domain
		}
	}()

	for {
		currentIP, err := godns.GetCurrentIP(handler.Configuration)
		if err != nil {
			log.Println("get_currentIP:", err)
			continue
		}
		log.Println("currentIP is:", currentIP)
		for _, subDomain := range domain.SubDomains {
			hostname := subDomain + "." + domain.DomainName
			lastIP := godns.ResolveDNS(hostname, handler.Configuration.Resolver, handler.Configuration.IPType)
			//check against currently known IP, if no change, skip update
			if currentIP == lastIP {
				log.Printf("IP is the same as cached one. Skip update.\n")
			} else {
				log.Printf("%s.%s Start to update record IP...\n", subDomain, domain.DomainName)
				handler.UpdateIP(domain.DomainName, subDomain, currentIP)

				// Send notification
				if err := godns.SendNotify(handler.Configuration, fmt.Sprintf("%s.%s", subDomain, domain.DomainName), currentIP); err != nil {
					log.Println("Failed to send notification")
				}
			}
		}
		// Sleep with interval
		log.Printf("Going to sleep, will start next checking in %d seconds...\r\n", handler.Configuration.Interval)
		time.Sleep(time.Second * time.Duration(handler.Configuration.Interval))
	}

}

// UpdateIP update subdomain with current IP
func (handler *Handler) UpdateIP(domain, subDomain, currentIP string) {
	client := godns.GetHttpClient(handler.Configuration)
	resp, err := client.Get(fmt.Sprintf(GoogleURL,
		handler.Configuration.Email,
		handler.Configuration.Password,
		subDomain,
		domain,
		currentIP))

	if err != nil {
		// handle error
		log.Print("Failed to update sub domain:", subDomain)
	}

	defer resp.Body.Close()

	if err != nil {
		log.Println("Request error...")
		log.Println("Err:", err.Error())
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode == http.StatusOK {
			if strings.Contains(string(body), "good") {
				log.Println("Update IP success:", string(body))
			} else if strings.Contains(string(body), "nochg") {
				log.Println("IP not changed:", string(body))
			}
		} else {
			log.Println("Update IP failed:", string(body))
		}
	}
}
