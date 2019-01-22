package dnspod

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/TimothyYe/godns"
	simplejson "github.com/bitly/go-simplejson"
	"golang.org/x/net/proxy"
)

// DNSPodHandler struct definition
type DNSPodHandler struct {
	Configuration *godns.Settings
}

// SetConfiguration pass dns settings and store it to handler instance
func (handler *DNSPodHandler) SetConfiguration(conf *godns.Settings) {
	handler.Configuration = conf
}

// DomainLoop the main logic loop
func (handler *DNSPodHandler) DomainLoop(domain *godns.Domain, panicChan chan<- godns.Domain) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Recovered in %v: %v\n", err, debug.Stack())
			panicChan <- *domain
		}
	}()

	var lastIP string
	for {
		log.Printf("Checking IP for domain %s \r\n", domain.DomainName)
		domainID := handler.GetDomain(domain.DomainName)

		if domainID == -1 {
			continue
		}

		currentIP, err := godns.GetCurrentIP(handler.Configuration)

		if err != nil {
			log.Println("get_currentIP:", err)
			continue
		}
		log.Println("currentIP is:", currentIP)

		//check against locally cached IP, if no change, skip update
		if currentIP == lastIP {
			log.Printf("IP is the same as cached one. Skip update.\n")
		} else {
			lastIP = currentIP

			for _, subDomain := range domain.SubDomains {

				subDomainID, ip := handler.GetSubDomain(domainID, subDomain)

				if subDomainID == "" || ip == "" {
					log.Printf("domain: %s.%s subDomainID: %s ip: %s\n", subDomain, domain.DomainName, subDomainID, ip)
					continue
				}

				// Continue to check the IP of subdomain
				if len(ip) > 0 && strings.TrimRight(currentIP, "\n") != strings.TrimRight(ip, "\n") {
					log.Printf("%s.%s Start to update record IP...\n", subDomain, domain.DomainName)
					handler.UpdateIP(domainID, subDomainID, subDomain, currentIP)

					// Send mail notification if notify is enabled
					if handler.Configuration.Notify.Enabled {
						log.Print("Sending notification to:", handler.Configuration.Notify.SendTo)
						godns.SendNotify(handler.Configuration, fmt.Sprintf("%s.%s", subDomain, domain.DomainName), currentIP)
					}

				} else {
					log.Printf("%s.%s Current IP is same as domain IP, no need to update...\n", subDomain, domain.DomainName)
				}
			}
		}
		// Interval is 5 minutes
		log.Printf("Going to sleep, will start next checking in %d minutes...\r\n", godns.INTERVAL)
		time.Sleep(time.Minute * godns.INTERVAL)
	}
}

// GenerateHeader generates the request header for DNSPod API
func (handler *DNSPodHandler) GenerateHeader(content url.Values) url.Values {
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
func (handler *DNSPodHandler) GetDomain(name string) int64 {

	var ret int64
	values := url.Values{}
	values.Add("type", "all")
	values.Add("offset", "0")
	values.Add("length", "20")

	response, err := handler.PostData("/Domain.List", values)

	if err != nil {
		log.Println("Failed to get domain list...")
		return -1
	}

	sjson, parseErr := simplejson.NewJson([]byte(response))

	if parseErr != nil {
		log.Println(parseErr)
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
			log.Println("domains slice is empty.")
		}
	} else {
		log.Println("get_domain:status code:", sjson.Get("status").Get("code").MustString())
	}

	return ret
}

// GetSubDomain returns subdomain by domain id
func (handler *DNSPodHandler) GetSubDomain(domainID int64, name string) (string, string) {
	log.Println("debug:", domainID, name)
	var ret, ip string
	value := url.Values{}
	value.Add("domain_id", strconv.FormatInt(domainID, 10))
	value.Add("offset", "0")
	value.Add("length", "1")
	value.Add("sub_domain", name)

	response, err := handler.PostData("/Record.List", value)

	if err != nil {
		log.Println("Failed to get domain list")
		return "", ""
	}

	sjson, parseErr := simplejson.NewJson([]byte(response))

	if parseErr != nil {
		log.Println(parseErr)
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
			log.Println("records slice is empty.")
		}
	} else {
		log.Println("get_subdomain:status code:", sjson.Get("status").Get("code").MustString())
	}

	return ret, ip
}

// UpdateIP update subdomain with current IP
func (handler *DNSPodHandler) UpdateIP(domainID int64, subDomainID string, subDomainName string, ip string) {
	value := url.Values{}
	value.Add("domain_id", strconv.FormatInt(domainID, 10))
	value.Add("record_id", subDomainID)
	value.Add("sub_domain", subDomainName)
	value.Add("record_type", "A")
	value.Add("record_line", "默认")
	value.Add("value", ip)

	response, err := handler.PostData("/Record.Modify", value)

	if err != nil {
		log.Println("Failed to update record to new IP!")
		log.Println(err)
		return
	}

	sjson, parseErr := simplejson.NewJson([]byte(response))

	if parseErr != nil {
		log.Println(parseErr)
		return
	}

	if sjson.Get("status").Get("code").MustString() == "1" {
		log.Println("New IP updated!")
	}

}

// PostData post data and invoke DNSPod API
func (handler *DNSPodHandler) PostData(url string, content url.Values) (string, error) {
	client := &http.Client{}

	if handler.Configuration.Socks5Proxy != "" {

		log.Println("use socks5 proxy:" + handler.Configuration.Socks5Proxy)

		dialer, err := proxy.SOCKS5("tcp", handler.Configuration.Socks5Proxy, nil, proxy.Direct)
		if err != nil {
			fmt.Println("can't connect to the proxy:", err)
			return "", err
		}

		httpTransport := &http.Transport{}
		client.Transport = httpTransport
		httpTransport.Dial = dialer.Dial
	}

	values := handler.GenerateHeader(content)
	req, _ := http.NewRequest("POST", "https://dnsapi.cn"+url, strings.NewReader(values.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", fmt.Sprintf("GoDNS/0.1 (%s)", ""))

	response, err := client.Do(req)

	if err != nil {
		log.Println("Post failed...")
		log.Println(err)
		return "", err
	}

	defer response.Body.Close()
	resp, _ := ioutil.ReadAll(response.Body)

	return string(resp), nil
}
