package alidns

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/TimothyYe/godns/internal/settings"
	"github.com/TimothyYe/godns/internal/utils"
	"github.com/TimothyYe/godns/pkg/notification"

	log "github.com/sirupsen/logrus"
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
	aliDNS := NewAliDNS(handler.Configuration.Email, handler.Configuration.Password, handler.Configuration.IPType)

	for while := true; while; while = !runOnce {
		if looping {
			// Sleep with interval
			log.Debugf("Going to sleep, will start next checking in %d seconds...", handler.Configuration.Interval)
			time.Sleep(time.Second * time.Duration(handler.Configuration.Interval))
		}

		looping = true
		currentIP, err := utils.GetCurrentIP(handler.Configuration)

		if err != nil {
			log.Error("Failed to get current IP:", err)
			continue
		}
		log.Debug("currentIP is:", currentIP)
		for _, subDomain := range domain.SubDomains {
			var hostname string
			if subDomain != utils.RootDomain {
				hostname = subDomain + "." + domain.DomainName
			} else {
				hostname = domain.DomainName
			}

			lastIP, err := utils.ResolveDNS(hostname, handler.Configuration.Resolver, handler.Configuration.IPType)
			if err != nil {
				log.Error(err)
				continue
			}
			//check against currently known IP, if no change, skip update
			if currentIP == lastIP {
				log.Infof("IP is the same as cached one (%s). Skip update.", currentIP)
			} else {
				log.Infof("%s.%s - Start to update record IP...", subDomain, domain.DomainName)
				records := aliDNS.GetDomainRecords(domain.DomainName, subDomain)
				if len(records) == 0 {
					log.Infof("Cannot get subdomain %s from AliDNS.", subDomain)
					continue
				}

				records[0].Value = currentIP
				if err := aliDNS.UpdateDomainRecord(records[0]); err != nil {
					log.Infof("Failed to update IP for subdomain:%s", subDomain)
					continue
				} else {
					log.Infof("IP updated for subdomain:%s", subDomain)
				}

				// Send notification
				notification.GetNotificationManager(handler.Configuration).
					Send(fmt.Sprintf("%s.%s", subDomain, domain.DomainName), currentIP)
			}
		}
	}

}
