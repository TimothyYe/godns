package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

var (
	HEUrl = "https://dyn.dns.he.net/nic/update"
)

type HEHandler struct{}

func (handler *HEHandler) DomainLoop(domain *Domain) {
	defer func() {
		if err := recover(); err != nil {
			panicCount++
			log.Printf("Recovered in %v: %v\n", err, debug.Stack())
			fmt.Println(identifyPanic())
			log.Print(identifyPanic())
			if panicCount < PANIC_MAX {
				log.Println("Got panic in goroutine, will start a new one... :", panicCount)
				go handler.DomainLoop(domain)
			} else {
				os.Exit(1)
			}
		}
	}()

	for {
		currentIP, err := getCurrentIP(configuration.IPUrl)

		if err != nil {
			log.Println("get_currentIP:", err)
			continue
		}
		log.Println("currentIp is:", currentIP)

		for _, subDomain := range domain.SubDomains {
			log.Printf("%s.%s Start to update record IP...\n", subDomain, domain.DomainName)
			handler.UpdateIP(domain.DomainName, subDomain, currentIP)
		}

		//Interval is 5 minutes
		time.Sleep(time.Minute * INTERVAL)
	}
}

func (handler *HEHandler) UpdateIP(domain, subDomain, currentIP string) {
	values := url.Values{}
	values.Add("hostname", fmt.Sprintf("%s.%s", subDomain, domain))
	values.Add("password", configuration.Password)
	values.Add("myip", currentIP)

	client := &http.Client{}

	if configuration.Socks5Proxy != "" {
		log.Println("use socks5 proxy:" + configuration.Socks5Proxy)
		dialer, err := proxy.SOCKS5("tcp", configuration.Socks5Proxy, nil, proxy.Direct)
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
