package duck

import (
	"fmt"
	"io/ioutil"
	"runtime/debug"
	"strings"
	"time"

	"github.com/TimothyYe/godns/internal/settings"
	"github.com/TimothyYe/godns/internal/utils"
	"github.com/TimothyYe/godns/pkg/notification"

	log "github.com/sirupsen/logrus"
)

var (
	// DuckURL the API address for Duck DNS.
	DuckURL = "https://www.duckdns.org/update?domains=%s&token=%s&%s"
)

// Handler struct.
type Handler struct {
	Configuration *settings.Settings
}

// SetConfiguration pass dns settings and store it to handler instance.
func (handler *Handler) SetConfiguration(conf *settings.Settings) {
	handler.Configuration = conf
}

// DomainLoop the main logic loop.
func (handler *Handler) DomainLoop(domain *settings.Domain, panicChan chan<- settings.Domain, runOnce bool) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("Recovered in %v: %v", err, string(debug.Stack()))
			panicChan <- *domain
		}
	}()

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
			log.Error("get_currentIP:", err)
			continue
		}

		log.Debug("currentIP is:", currentIP)
		client := utils.GetHTTPClient(handler.Configuration, handler.Configuration.UseProxy)
		var ip string

		if strings.ToUpper(handler.Configuration.IPType) == utils.IPV4 {
			ip = fmt.Sprintf("ip=%s", currentIP)
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

			//check against currently known IP, if no change, skip update
			if currentIP == lastIP {
				log.Infof("IP is the same as cached one (%s). Skip update.", currentIP)
			} else {
				// update IP with HTTP GET request
				resp, err := client.Get(fmt.Sprintf(DuckURL, subDomain, handler.Configuration.LoginToken, ip))
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
					log.Infof("IP updated to: %s", currentIP)
				}

				// Send notification
				notification.GetNotificationManager(handler.Configuration).Send(fmt.Sprintf("%s.%s", subDomain, domain.DomainName), currentIP)
			}
		}
	}
}
