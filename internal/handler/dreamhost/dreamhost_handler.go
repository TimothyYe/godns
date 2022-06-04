package dreamhost

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

	"github.com/google/uuid"
)

var (
	// DreamhostURL the API address for dreamhost.com.
	DreamhostURL = "https://api.dreamhost.com"
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
				log.Println(err)
				continue
			}

			//check against currently known IP, if no change, skip update
			if currentIP == lastIP {
				log.Infof("IP is the same as cached one (%s). Skip update.", currentIP)
			} else {
				log.Infof("%s.%s Start to update record IP...", subDomain, domain.DomainName)
				handler.UpdateIP(hostname, currentIP, lastIP)

				// Send notification
				notification.GetNotificationManager(handler.Configuration).Send(fmt.Sprintf("%s.%s", subDomain, domain.DomainName), currentIP)
			}
		}
	}

}

// UpdateIP update subdomain with current IP.
func (handler *Handler) UpdateIP(hostname, currentIP, lastIP string) {

	handler.updateDNS(lastIP, currentIP, hostname, "remove")
	handler.updateDNS(lastIP, currentIP, hostname, "add")

}

// updateDNS can add or remove DNS records.
func (handler *Handler) updateDNS(dns, ip, hostname, action string) {
	var ipType string
	if handler.Configuration.IPType == "" || strings.ToUpper(handler.Configuration.IPType) == utils.IPV4 {
		ipType = utils.IPTypeA
	} else if strings.ToUpper(handler.Configuration.IPType) == utils.IPV6 {
		ipType = utils.IPTypeAAAA
	}

	// Generates UUID
	uid, _ := uuid.NewRandom()
	values := url.Values{}
	values.Add("record", hostname)
	values.Add("key", handler.Configuration.LoginToken)
	values.Add("type", ipType)
	values.Add("unique_id", uid.String())
	switch action {
	case "remove":
		// Build URL query (remove)
		values.Add("cmd", "dns-remove_record")
		values.Add("value", dns)
	case "add":
		// Build URL query (add)
		values.Add("cmd", "dns-add_record")
		values.Add("value", ip)
	default:
		log.Fatalf("Unknown action %s", action)
	}

	client := utils.GetHTTPClient(handler.Configuration, handler.Configuration.UseProxy)
	req, _ := http.NewRequest("POST", DreamhostURL, strings.NewReader(values.Encode()))
	req.SetBasicAuth(handler.Configuration.Email, handler.Configuration.Password)

	if handler.Configuration.UserAgent != "" {
		req.Header.Add("User-Agent", handler.Configuration.UserAgent)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Error("Request err:", err.Error())
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode == http.StatusOK {
			log.Infof("Update IP success: %s", string(body))
		} else {
			log.Infof("Update IP failed: %s", string(body))
		}
	}
}
