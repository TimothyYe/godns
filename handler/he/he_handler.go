package he

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"runtime/debug"
	"strings"
	"time"

	"github.com/TimothyYe/godns"
)

var (
	// HEUrl the API address for he.net
	HEUrl = "https://dyn.dns.he.net/nic/update"
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

	var lastIP string
	for {
		currentIP, err := godns.GetCurrentIP(handler.Configuration)

		if err != nil {
			log.Println("get_currentIP:", err)
			continue
		}
		log.Println("currentIP is:", currentIP)

		//check against locally cached IP, if no change, skip update
		if currentIP == lastIP {
			log.Printf("IP is the same as cached one. Skip update.\n")
		} else {
			lastIP = currentIP

			for _, subDomain := range domain.SubDomains {
				log.Printf("%s.%s Start to update record IP...\n", subDomain, domain.DomainName)
				handler.UpdateIP(domain.DomainName, subDomain, currentIP)

				// Send mail notification if notify is enabled
				if handler.Configuration.Notify.Enabled {
					log.Print("Sending notification to:", handler.Configuration.Notify.SendTo)
					if err := godns.SendNotify(handler.Configuration, fmt.Sprintf("%s.%s", subDomain, domain.DomainName), currentIP); err != nil {
						log.Println("Failed to send notificaiton")
					}
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
	values := url.Values{}
	values.Add("hostname", fmt.Sprintf("%s.%s", subDomain, domain))
	values.Add("password", handler.Configuration.Password)
	values.Add("myip", currentIP)

	client := godns.GetHttpClient(handler.Configuration)

	req, _ := http.NewRequest("POST", HEUrl, strings.NewReader(values.Encode()))
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
