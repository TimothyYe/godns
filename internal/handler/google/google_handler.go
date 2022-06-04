package google

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/TimothyYe/godns/internal/settings"
	"github.com/TimothyYe/godns/internal/utils"
	"github.com/TimothyYe/godns/pkg/notification"

	log "github.com/sirupsen/logrus"
)

var (
	// GoogleURL the API address for Google Domains.
	GoogleURL = "https://%s:%s@domains.google.com/nic/update?hostname=%s.%s&myip=%s"
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
				log.Infof("%s.%s Start to update record IP...", subDomain, domain.DomainName)
				handler.UpdateIP(domain.DomainName, subDomain, currentIP)

				// Send notification
				notification.GetNotificationManager(handler.Configuration).Send(fmt.Sprintf("%s.%s", subDomain, domain.DomainName), currentIP)
			}
		}
	}

}

// UpdateIP update subdomain with current IP.
func (handler *Handler) UpdateIP(domain, subDomain, currentIP string) {
	client := utils.GetHTTPClient(handler.Configuration, handler.Configuration.UseProxy)
	resp, err := client.Get(fmt.Sprintf(GoogleURL,
		handler.Configuration.Email,
		handler.Configuration.Password,
		subDomain,
		domain,
		currentIP))

	if err != nil {
		// handle error
		log.Error("Failed to update sub domain:", subDomain)
		return
	}

	defer resp.Body.Close()

	if err != nil {
		log.Error("Err:", err.Error())
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode == http.StatusOK {
			if strings.Contains(string(body), "good") {
				log.Infof("Update IP success: %s", string(body))
			} else if strings.Contains(string(body), "nochg") {
				log.Infof("IP not changed: %s", string(body))
			}
		} else {
			log.Infof("Update IP failed: %s", string(body))
		}
	}
}
