package handler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"runtime/debug"
	"strings"
	"time"

	"github.com/TimothyYe/godns"

	"golang.org/x/net/proxy"
)

var (
	// HEUrl the API address for he.net
	HEUrl = "https://dyn.dns.he.net/nic/update"
)

// HEHandler struct
type HEHandler struct {
	Configuration *godns.Settings
}

func (handler *HEHandler) SetConfiguration(conf *godns.Settings) {
	handler.Configuration = conf
}

// DomainLoop the main logic loop
func (handler *HEHandler) DomainLoop(domain *godns.Domain, panicChan chan<- godns.Domain) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Recovered in %v: %v\n", err, debug.Stack())
			fmt.Println(godns.IdentifyPanic())
			log.Print(godns.IdentifyPanic())
			panicChan <- *domain
		}
	}()

	for {
		currentIP, err := godns.GetCurrentIP(handler.Configuration)

		if err != nil {
			log.Println("get_currentIP:", err)
			continue
		}
		log.Println("currentIp is:", currentIP)

		for _, subDomain := range domain.SubDomains {
			log.Printf("%s.%s Start to update record IP...\n", subDomain, domain.DomainName)
			handler.UpdateIP(domain.DomainName, subDomain, currentIP)
		}

		// Interval is 5 minutes
		time.Sleep(time.Minute * godns.INTERVAL)
	}
}

// UpdateIP update subdomain with current IP
func (handler *HEHandler) UpdateIP(domain, subDomain, currentIP string) {
	values := url.Values{}
	values.Add("hostname", fmt.Sprintf("%s.%s", subDomain, domain))
	values.Add("password", handler.Configuration.Password)
	values.Add("myip", currentIP)

	client := &http.Client{}

	if handler.Configuration.Socks5Proxy != "" {
		log.Println("use socks5 proxy:" + handler.Configuration.Socks5Proxy)
		dialer, err := proxy.SOCKS5("tcp", handler.Configuration.Socks5Proxy, nil, proxy.Direct)
		if err != nil {
			log.Println("can't connect to the proxy:", err)
			return
		}

		httpTransport := &http.Transport{}
		client.Transport = httpTransport
		httpTransport.Dial = dialer.Dial
	}

	req, _ := http.NewRequest("POST", HEUrl, strings.NewReader(values.Encode()))
	resp, err := client.Do(req)

	if err != nil {
		log.Println("Request error...")
		log.Println("Err:", err.Error())
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode == http.StatusOK {
			log.Println("Update IP success:", string(body))
		} else {
			log.Println("Update IP failed:", string(body))
		}
	}
}
