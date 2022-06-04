package dynv6

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
	// Dynv6URL the API address for Duck DNS.
	Dynv6URL = "https://dynv6.com/api/update?hostname=%s&token=%s&%s"
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
				log.Infof("IP is the same as cached one (%s). Skip update.", currentIP)
			} else {
				err := handler.update(client, hostname, currentIP)
				if err != nil {
					log.Errorf("Failed to update domain %s: %v", hostname, err)
				} else {
					log.Infof("IP updated to: %s", currentIP)
				}
			}
		}
	}
}

func (handler *Handler) update(client *http.Client, hostname string, currentIP string) error {

	var ip string
	if strings.ToUpper(handler.Configuration.IPType) == utils.IPV4 {
		ip = fmt.Sprintf("ipv4=%s", currentIP)
	} else if strings.ToUpper(handler.Configuration.IPType) == utils.IPV6 {
		ip = fmt.Sprintf("ipv6=%s", currentIP)
	}

	// update IP with HTTP GET request
	url := fmt.Sprintf(Dynv6URL, hostname, handler.Configuration.LoginToken, ip)
	log.Debug("Update url: ", url)
	resp, err := client.Get(url)
	if err != nil {
		// handle error
		return fmt.Errorf("cannot send request: %w", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to receive response: %w", err)
	} else if !strings.Contains(string(body), "addresses updated") {
		return fmt.Errorf("service rejected update: %s", string(body))
	}

	// Send notification
	notification.GetNotificationManager(handler.Configuration).Send(hostname, currentIP)
	return nil
}
