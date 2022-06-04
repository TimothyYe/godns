package scaleway

import (
	"bytes"
	"encoding/json"
	"errors"
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
	ScalewayURL = "https://api.scaleway.com/domain/v2beta1/dns-zones/%s/records"
)

// Handler struct.
type Handler struct {
	Configuration *settings.Settings
}

// Record for Scaleway API.
type Record struct {
	Name    string `json:"name"`
	Data    string `json:"data"`
	TTL     int    `json:"ttl"`
	Comment string `json:"comment"`
}

// IDFields to filter DNS records for Scaleway API.
type IDFields struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// SetRecord for Scaleway API.
type SetRecord struct {
	IDFields IDFields `json:"id_fields"`
	Records  []Record `json:"records"`
}

// DNSChange for Scaleway API.
type DNSChange struct {
	Set SetRecord `json:"set"`
}

// DNSUpdateRequest for Scaleway API.
type DNSUpdateRequest struct {
	Changes []DNSChange `json:"changes"`
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
			var hostname string
			if subDomain != utils.RootDomain {
				hostname = subDomain + "." + domain.DomainName
			} else {
				hostname = domain.DomainName
			}
			lastIP, err := utils.ResolveDNS(hostname, handler.Configuration.Resolver, handler.Configuration.IPType)
			if err != nil {
				log.Error(err)
				recordType, _ := handler.GetRecordType()
				log.Errorf("Failed to resolve '%s', make sure a '%s' record exists.", hostname, recordType)
				continue
			}

			// Check against currently known IP, if no change, skip update
			if currentIP == lastIP {
				log.Infof("IP is the same as cached one (%s). Skip update.", currentIP)
			} else {
				log.Infof("%s.%s - Start to update record IP...", subDomain, domain.DomainName)
				err := handler.UpdateIP(domain.DomainName, subDomain, currentIP)
				if err != nil {
					log.Error(err)
					continue
				}
				// Send notification
				notification.GetNotificationManager(handler.Configuration).Send(fmt.Sprintf("%s.%s", subDomain, domain.DomainName), currentIP)
			}
		}
	}
}

func (handler *Handler) GetRecordType() (string, error) {
	if strings.ToUpper(handler.Configuration.IPType) == utils.IPV4 {
		return utils.IPTypeA, nil
	} else if strings.ToUpper(handler.Configuration.IPType) == utils.IPV6 {
		return utils.IPTypeAAAA, nil
	} else {
		return "", errors.New("must specify \"ip_type\" in config for Scaleway")
	}
}

// UpdateIP update subdomain with current IP.
func (handler *Handler) UpdateIP(domain, subDomain, currentIP string) error {
	recordType, err := handler.GetRecordType()
	if err != nil {
		return err
	}

	reqBody := DNSUpdateRequest{Changes: []DNSChange{{SetRecord{
		IDFields: IDFields{
			Name: subDomain,
			Type: recordType,
		},
		Records: []Record{
			{
				Name:    subDomain,
				Data:    currentIP,
				TTL:     handler.Configuration.Interval,
				Comment: "Set by godns",
			},
		},
	}}}}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return errors.New("failed to encode request body as json")
	}

	req, _ := http.NewRequest("PATCH", fmt.Sprintf(ScalewayURL, domain), bytes.NewReader(jsonBody))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Auth-Token", handler.Configuration.LoginToken)
	if handler.Configuration.UserAgent != "" {
		req.Header.Add("User-Agent", handler.Configuration.UserAgent)
	}

	client := utils.GetHTTPClient(handler.Configuration, handler.Configuration.UseProxy)
	log.Debugf("Requesting update for '%s.%s': '%v'", subDomain, domain, reqBody)
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Request error:", err)
		return errors.New("failed to complete update request")
	}

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Debugf("Update failed for '%s.%s': %s", subDomain, domain, string(body))
		return fmt.Errorf("update IP failed with status '%d'", resp.StatusCode)
	}
	log.Debugf("Update IP success for '%s.%s': '%s'", subDomain, domain, string(body))
	return nil
}
