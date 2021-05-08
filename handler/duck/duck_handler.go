package duck

import (
	"fmt"
	"io/ioutil"
	"runtime/debug"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/TimothyYe/godns/notify"

	"github.com/TimothyYe/godns"
)

var (
	// DuckUrl the API address for Duck DNS
	DuckUrl = "https://www.duckdns.org/update?domains=%s&token=%s&%s"
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
			log.Errorf("Recovered in %v: %v\n", err, string(debug.Stack()))
			panicChan <- *domain
		}
	}()

	looping := false

	for {
		if looping {
			// Sleep with interval
			log.Debug("Going to sleep, will start next checking in %d seconds...\r\n", handler.Configuration.Interval)
			time.Sleep(time.Second * time.Duration(handler.Configuration.Interval))
		}

		looping = true
		currentIP, err := godns.GetCurrentIP(handler.Configuration)

		if err != nil {
			log.Error("get_currentIP:", err)
			continue
		}

		log.Debug("currentIP is:", currentIP)
		client := godns.GetHttpClient(handler.Configuration, handler.Configuration.UseProxy)
		var ip string

		if strings.ToUpper(handler.Configuration.IPType) == godns.IPV4 {
			ip = fmt.Sprintf("ip=%s", currentIP)
		} else if strings.ToUpper(handler.Configuration.IPType) == godns.IPV6 {
			ip = fmt.Sprintf("ipv6=%s", currentIP)
		}

		for _, subDomain := range domain.SubDomains {
			hostname := subDomain + "." + domain.DomainName
			lastIP, err := godns.ResolveDNS(hostname, handler.Configuration.Resolver, handler.Configuration.IPType)
			if err != nil {
				log.Error(err)
				continue
			}

			//check against currently known IP, if no change, skip update
			if currentIP == lastIP {
				log.Infof("IP is the same as cached one (%s). Skip update.\n", currentIP)
			} else {
				// update IP with HTTP GET request
				resp, err := client.Get(fmt.Sprintf(DuckUrl, subDomain, handler.Configuration.LoginToken, ip))
				if err != nil {
					// handle error
					log.Error("Failed to update sub domain:", subDomain)
					continue
				}

				defer resp.Body.Close()

				body, err := ioutil.ReadAll(resp.Body)
				if err != nil || string(body) != "OK" {
					log.Error("Failed to update the IP:", err)
					continue
				} else {
					log.Info("IP updated to:", currentIP)
				}

				// Send notification
				notify.GetNotifyManager(handler.Configuration).Send(fmt.Sprintf("%s.%s", subDomain, domain.DomainName), currentIP)
			}
		}
	}
}
