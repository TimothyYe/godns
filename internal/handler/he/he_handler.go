package he

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime/debug"
	"strings"
	"time"

	"github.com/TimothyYe/godns/internal/settings"
	"github.com/TimothyYe/godns/internal/utils"
	"github.com/TimothyYe/godns/pkg/notification"

	log "github.com/sirupsen/logrus"
)

var (
	// HEUrl the API address for he.net.
	HEUrl = "https://dyn.dns.he.net/nic/update"
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

		//check against locally cached IP, if no change, skip update

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
				handler.UpdateIP(domain.DomainName, subDomain, currentIP)

				// Send notification
				notification.GetNotificationManager(handler.Configuration).Send(fmt.Sprintf("%s.%s", subDomain, domain.DomainName), currentIP)
			}
		}
	}

}

// UpdateIP update subdomain with current IP.
func (handler *Handler) UpdateIP(domain, subDomain, currentIP string) {
	values := url.Values{}

	if subDomain != utils.RootDomain {
		values.Add("hostname", fmt.Sprintf("%s.%s", subDomain, domain))
	} else {
		values.Add("hostname", domain)
	}
	values.Add("password", handler.Configuration.Password)
	values.Add("myip", currentIP)

	client := utils.GetHTTPClient(handler.Configuration, handler.Configuration.UseProxy)

	req, _ := http.NewRequest("POST", HEUrl, strings.NewReader(values.Encode()))
	resp, err := client.Do(req)

	if err != nil {
		log.Error("Request error:", err)
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode == http.StatusOK {
			log.Infof("Update IP success: %s", string(body))
		} else {
			log.Infof("Update IP failed: %s", string(body))
		}
	}
}
