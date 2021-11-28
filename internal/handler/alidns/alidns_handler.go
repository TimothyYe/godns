package alidns

import (
	"fmt"
	"github.com/TimothyYe/godns/internal/settings"
	"github.com/TimothyYe/godns/internal/utils"
	"github.com/TimothyYe/godns/pkg/notify"
	"runtime/debug"
	"time"

	log "github.com/sirupsen/logrus"
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
	aliDNS := NewAliDNS(handler.Configuration.Email, handler.Configuration.Password, handler.Configuration.IPType)

	for {
		if looping {
			// Sleep with interval
			log.Debugf("Going to sleep, will start next checking in %d seconds...\r\n", handler.Configuration.Interval)
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
				log.Infof("IP is the same as cached one (%s). Skip update.\n", currentIP)
			} else {
				lastIP = currentIP

				log.Infof("%s.%s - Start to update record IP...\n", subDomain, domain.DomainName)
				records := aliDNS.GetDomainRecords(domain.DomainName, subDomain)
				if records == nil || len(records) == 0 {
					log.Infof("Cannot get subdomain %s from AliDNS.\r\n", subDomain)
					continue
				}

				records[0].Value = currentIP
				if err := aliDNS.UpdateDomainRecord(records[0]); err != nil {
					log.Infof("Failed to update IP for subdomain:%s\r\n", subDomain)
					continue
				} else {
					log.Infof("IP updated for subdomain:%s\r\n", subDomain)
				}

				// Send notification
				notify.GetNotifyManager(handler.Configuration).
					Send(fmt.Sprintf("%s.%s", subDomain, domain.DomainName), currentIP)
			}
		}
	}

}
