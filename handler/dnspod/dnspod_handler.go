package dnspod

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/TimothyYe/godns/notify"

	"github.com/TimothyYe/godns"
	"github.com/bitly/go-simplejson"
)

// Handler struct definition
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
			log.Debugf("Going to sleep, will start next checking in %d seconds...\r\n", handler.Configuration.Interval)
			time.Sleep(time.Second * time.Duration(handler.Configuration.Interval))
		}

		looping = true

		log.Infof("Checking IP for domain %s \r\n", domain.DomainName)
		domainID := handler.GetDomain(domain.DomainName)

		if domainID == -1 {
			continue
		}

		currentIP, err := godns.GetCurrentIP(handler.Configuration)

		if err != nil {
			log.Error("get_currentIP:", err)
			continue
		}
		log.Debug("currentIP is:", currentIP)

		for _, subDomain := range domain.SubDomains {
			var hostname string
			if subDomain != godns.RootDomain {
				hostname = subDomain + "." + domain.DomainName
			} else {
				hostname = domain.DomainName
			}

			lastIP, err := godns.ResolveDNS(hostname, handler.Configuration.Resolver, handler.Configuration.IPType)
			if err != nil {
				log.Println(err)
				continue
			}

			//check against currently known IP, if no change, skip update
			if currentIP == lastIP {
				log.Infof("IP is the same as cached one (%s). Skip update.\n", currentIP)
			} else {
				lastIP = currentIP

				subDomainID, ip := handler.GetSubDomain(domainID, subDomain)

				if subDomainID == "" || ip == "" {
					log.Infof("Domain or subdomain not configured yet. domain: %s.%s subDomainID: %s ip: %s\n", subDomain, domain.DomainName, subDomainID, ip)
					continue
				}

				// Continue to check the IP of subdomain
				if len(ip) > 0 && strings.TrimRight(currentIP, "\n") != strings.TrimRight(ip, "\n") {
					log.Infof("%s.%s Start to update record IP...\n", subDomain, domain.DomainName)
					handler.UpdateIP(domainID, subDomainID, subDomain, currentIP)

					// Send notification
					notify.GetNotifyManager(handler.Configuration).Send(fmt.Sprintf("%s.%s", subDomain, domain.DomainName), currentIP)
				} else {
					log.Infof("%s.%s Current IP is same as domain IP, no need to update...\n", subDomain, domain.DomainName)
				}
			}
		}
	}
}

// GenerateHeader generates the request header for DNSPod API
func (handler *Handler) GenerateHeader(content url.Values) url.Values {
	header := url.Values{}
	if handler.Configuration.LoginToken != "" {
		header.Add("login_token", handler.Configuration.LoginToken)
	}

	header.Add("format", "json")
	header.Add("lang", "en")
	header.Add("error_on_empty", "no")

	if content != nil {
		for k := range content {
			header.Add(k, content.Get(k))
		}
	}

	return header
}

// GetDomain returns specific domain by name
func (handler *Handler) GetDomain(name string) int64 {

	var ret int64
	values := url.Values{}
	values.Add("type", "all")
	values.Add("offset", "0")
	values.Add("length", "20")

	response, err := handler.PostData("/Domain.List", values)

	if err != nil {
		log.Error("Failed to get domain list:", err)
		return -1
	}

	sjson, parseErr := simplejson.NewJson([]byte(response))

	if parseErr != nil {
		log.Error(parseErr)
		return -1
	}

	if sjson.Get("status").Get("code").MustString() == "1" {
		domains, _ := sjson.Get("domains").Array()

		for _, d := range domains {
			m := d.(map[string]interface{})
			if m["name"] == name {
				id := m["id"]

				switch t := id.(type) {
				case json.Number:
					ret, _ = t.Int64()
				}

				break
			}
		}
		if len(domains) == 0 {
			log.Info("domains slice is empty.")
		}
	} else {
		log.Info("get_domain:status code:", sjson.Get("status").Get("code").MustString())
	}

	return ret
}

// GetSubDomain returns subdomain by domain id
func (handler *Handler) GetSubDomain(domainID int64, name string) (string, string) {
	var ret, ip string
	value := url.Values{}
	value.Add("domain_id", strconv.FormatInt(domainID, 10))
	value.Add("offset", "0")
	value.Add("length", "1")
	value.Add("sub_domain", name)

	if handler.Configuration.IPType == "" || strings.ToUpper(handler.Configuration.IPType) == godns.IPV4 {
		value.Add("record_type", "A")
	} else if strings.ToUpper(handler.Configuration.IPType) == godns.IPV6 {
		value.Add("record_type", "AAAA")
	} else {
		log.Error("Error: must specify \"ip_type\" in config for DNSPod.")
		return "", ""
	}

	response, err := handler.PostData("/Record.List", value)

	if err != nil {
		log.Error("Failed to get domain list:", err)
		return "", ""
	}

	sjson, parseErr := simplejson.NewJson([]byte(response))

	if parseErr != nil {
		log.Error(parseErr)
		return "", ""
	}

	if sjson.Get("status").Get("code").MustString() == "1" {
		records, _ := sjson.Get("records").Array()

		for _, d := range records {
			m := d.(map[string]interface{})
			if m["name"] == name {
				ret = m["id"].(string)
				ip = m["value"].(string)
				break
			}
		}
		if len(records) == 0 {
			log.Info("records slice is empty.")
		}
	} else {
		log.Info("get_subdomain:status code:", sjson.Get("status").Get("code").MustString())
	}

	return ret, ip
}

// UpdateIP update subdomain with current IP
func (handler *Handler) UpdateIP(domainID int64, subDomainID string, subDomainName string, ip string) {
	value := url.Values{}
	value.Add("domain_id", strconv.FormatInt(domainID, 10))
	value.Add("record_id", subDomainID)
	value.Add("sub_domain", subDomainName)

	if strings.ToUpper(handler.Configuration.IPType) == godns.IPV4 {
		value.Add("record_type", godns.IPTypeA)
	} else if strings.ToUpper(handler.Configuration.IPType) == godns.IPV6 {
		value.Add("record_type", godns.IPTypeAAAA)
	} else {
		log.Info("Error: must specify \"ip_type\" in config for DNSPod.")
		return
	}

	value.Add("record_line", "默认")
	value.Add("value", ip)

	response, err := handler.PostData("/Record.Modify", value)

	if err != nil {
		log.Error("Failed to update record to new IP:", err)
		return
	}

	sjson, parseErr := simplejson.NewJson([]byte(response))

	if parseErr != nil {
		log.Error(parseErr)
		return
	}

	if sjson.Get("status").Get("code").MustString() == "1" {
		log.Info("New IP updated!")
	} else {
		log.Info("Failed to update IP record:", sjson.Get("status").Get("message").MustString())
	}

}

// PostData post data and invoke DNSPod API
func (handler *Handler) PostData(url string, content url.Values) (string, error) {
	client := godns.GetHttpClient(handler.Configuration, handler.Configuration.UseProxy)

	if client == nil {
		return "", errors.New("failed to create HTTP client")
	}

	values := handler.GenerateHeader(content)
	req, _ := http.NewRequest("POST", "https://dnsapi.cn"+url, strings.NewReader(values.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", fmt.Sprintf("GoDNS/0.1 (%s)", ""))

	response, err := client.Do(req)

	if err != nil {
		log.Error("Post failed:", err)
		return "", err
	}

	defer response.Body.Close()
	resp, _ := ioutil.ReadAll(response.Body)

	return string(resp), nil
}
