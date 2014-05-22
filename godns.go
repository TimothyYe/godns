package main

import (
	"fmt"
	log "github.com/cihub/seelog"
	"os"
	"strings"
	"time"
)

var Configuration Settings

func main() {
	initLog()

	fmt.Println("Starting...")
	log.Info("Starting...")

	Configuration = LoadSettings()

	loop := make(chan bool)
	go dns_loop(loop)

	ret := <-loop

	if !ret {
		fmt.Println("Dns loop exited...")
		log.Error("Dns loop exited...")
		close(loop)

		os.Exit(1)
	}
}

func dns_loop(loop chan bool) {

	for {

		defer func() {
			if err := recover(); err != nil {
				log.Error(err)
			}
		}()

		domain_id := get_domain(Configuration.Domain)

		if domain_id == -1 {
			continue
		}

		currentIP, err := get_currentIP(Configuration.IP_Url)

		if err != nil {
			continue
		}

		sub_domain_id, ip := get_subdomain(domain_id, Configuration.Sub_domain)

		if sub_domain_id == "" || ip == "" {
			continue
		}

		fmt.Printf("currentIp is:%s\n", currentIP)
		log.Infof("currentIp is:%s\n", currentIP)

		//Continue to check the IP of sub-domain
		if len(ip) > 0 && !strings.Contains(currentIP, ip) {

			fmt.Println("Start to update record IP...")
			log.Info("Start to update record IP...")
			update_ip(domain_id, sub_domain_id, Configuration.Sub_domain, currentIP)

		} else {
			fmt.Println("Current IP is same as domain IP, no need to update...")
			log.Info("Current IP is same as domain IP, no need to update...")
		}

		//Interval is 5 minutes
		time.Sleep(time.Second * 60 * 5)
	}

	log.Info("Loop exited...")
	loop <- false
}
