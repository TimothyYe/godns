package dynv6

import (
	"fmt"
	"io/ioutil"
	"runtime/debug"
	"strings"
	"time"

	"github.com/TimothyYe/godns/internal/settings"
	"github.com/TimothyYe/godns/internal/utils"
	"github.com/TimothyYe/godns/pkg/notify"

	log "github.com/sirupsen/logrus"
)

var (
	// Dynv6Url the API address for Duck DNS
	Dynv6Url = "https://dynv6.com/api/update?hostname=%s&token=%s&%s"
)

// Handler struct
type Handler struct {
	Configuration *settings.Settings
}

// SetConfiguration pass dns settings and store it to handler instance
func (handler *Handler) SetConfiguration(conf *settings.Settings) {
	handler.Configuration = conf
}

// DomainLoop the main logic loop
func (handler *Handler) DomainLoop(domain *settings.Domain, panicChan chan<- settings.Domain) {
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
			log.Debugf("Going to sleep, will start next checking in %d seconds...\r\n", handler.Configuration.Interval)
			time.Sleep(time.Second * time.Duration(handler.Configuration.Interval))
		}

		looping = true
		currentIP, err := utils.GetCurrentIP(handler.Configuration)

		if err != nil {
			log.Error("get_currentIP:", err)
			continue
		}

		log.Debug("currentIP is:", currentIP)
		client := utils.GetHttpClient(handler.Configuration, handler.Configuration.UseProxy)
		var ip string

		if strings.ToUpper(handler.Configuration.IPType) == utils.IPV4 {
			ip = fmt.Sprintf("ipv4=%s", currentIP)
		} else if strings.ToUpper(handler.Configuration.IPType) == utils.IPV6 {
			ip = fmt.Sprintf("ipv6=%s", currentIP)
		}

		for _, subDomain := range domain.SubDomains {
			hostname := subDomain + "." + domain.DomainName
			lastIP, err := utils.ResolveDNS(hostname, handler.Configuration.Resolver, handler.Configuration.IPType)
			if err != nil {
				log.Error(err)
				continue
			}
			log.Debug("DNS record IP is:", lastIP)

			//check against currently known IP, if no change, skip update
			if currentIP == lastIP {
				log.Infof("IP is the same as cached one (%s). Skip update.\n", currentIP)
			} else {
				// update IP with HTTP GET request
				url := fmt.Sprintf(Dynv6Url, hostname, handler.Configuration.LoginToken, ip)
				log.Debug("Update url: ", url)
				resp, err := client.Get(url)
				if err != nil {
					// handle error
					log.Errorf("Failed to update domain %s: %v", hostname, err)
					continue
				}

				defer resp.Body.Close()

				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Error("Failed to update the IP: ", err)
					continue
				} else if string(body) != "OK" {
					log.Error("Failed to update the IP: ", string(body))
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
