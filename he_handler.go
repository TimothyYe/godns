package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/parnurzeal/gorequest"
)

var (
	HEUrl = "https://dyn.dns.he.net/nic/update?hostname=%s.%s&password=%s&myip=%s"
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
	request := gorequest.New()
	resp, body, errs := request.Get(fmt.Sprintf(HEUrl, subDomain, domain, configuration.Password, currentIP)).End()

	if len(errs) > 0 {
		log.Println("Request error...")

		for _, err := range errs {
			log.Println("Err:", err.Error())
		}
	} else {
		if resp.StatusCode == http.StatusOK {
			log.Println("Update IP success:", body)
		} else {
			log.Println("Update IP failed:", body)
		}
	}
}
