package noip

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/TimothyYe/godns/notify"

	"github.com/TimothyYe/godns"
)

var (
	// NoIPUrl the API address for NoIP
	NoIPUrl = "https://%s:%s@dynupdate.no-ip.com/nic/update?hostname=%s&%s"
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
		client := godns.GetHttpClient(handler.Configuration, handler.Configuration.UseProxy)

		var ip string
		if strings.ToUpper(handler.Configuration.IPType) == godns.IPV4 {
			ip = fmt.Sprintf("myip=%s", currentIP)
		} else if strings.ToUpper(handler.Configuration.IPType) == godns.IPV6 {
			ip = fmt.Sprintf("myipv6=%s", currentIP)
		}

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
				req, _ := http.NewRequest("GET", fmt.Sprintf(
					NoIPUrl,
					handler.Configuration.Email,
					handler.Configuration.Password,
					hostname,
					ip), nil)

				if handler.Configuration.UserAgent != "" {
					req.Header.Add("User-Agent", handler.Configuration.UserAgent)
				}

				// update IP with HTTP GET request
				resp, err := client.Do(req)
				if err != nil {
					// handle error
					log.Print("Failed to update sub domain:", subDomain)
					continue
				}

				defer resp.Body.Close()

				body, err := ioutil.ReadAll(resp.Body)
				if err != nil || !strings.Contains(string(body), "good") {
					log.Println("Failed to update the IP")
					continue
				} else {
					log.Print("IP updated to:", currentIP)
				}

				// Send notification
				notify.GetNotifyManager(handler.Configuration).Send(fmt.Sprintf("%s.%s", subDomain, domain.DomainName), currentIP)
			}
		}
	}
}
