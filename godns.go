package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var Configuration Settings

func main() {
	fmt.Println("Starting...")

	Configuration = LoadSettings()

	loop := make(chan bool)
	go dns_loop(loop)

	ret := <-loop

	if !ret {
		fmt.Println("Dns loop exited...")
		close(loop)

		os.Exit(1)
	}
}

func dns_loop(loop chan bool) {

	for {

		domain_id := get_domain(Configuration.Domain)

		currentIP, _ := get_currentIP(Configuration.IP_Url)
		sub_domain_id, ip := get_subdomain(domain_id, Configuration.Sub_domain)

		fmt.Printf("currentIp is:%s\n", currentIP)

		//Continue to check the IP of sub-domain
		if len(ip) > 0 && !strings.Contains(currentIP, ip) {

			fmt.Println("Start to update record IP...")
			update_ip(domain_id, sub_domain_id, Configuration.Sub_domain, currentIP)

		} else {
			fmt.Println("Current IP is same as domain IP, no need to update...")
		}

		//Interval is 5 minutes
		time.Sleep(time.Second * 60 * 5)
	}
}
