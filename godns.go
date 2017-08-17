package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

const (
	//PANIC_MAX is the max allowed panic times
	PANIC_MAX = 5
	//INTERVAL is minute
	INTERVAL = 5
)

var (
	configuration Settings
	optConf       = flag.String("c", "./config.json", "Specify a config file")
	optDocker     = flag.Bool("d", false, "Run it as docker mode")
	optHelp       = flag.Bool("h", false, "Show help")
	panicCount    = 0
)

func main() {
	flag.Parse()
	if *optHelp {
		flag.Usage()
		return
	}

	if *optDocker {
		//Load settings from ENV
		configuration = Settings{
			Email:      os.Getenv("EMAIL"),
			Password:   os.Getenv("PASSWORD"),
			LoginToken: os.Getenv("TOKEN"),
			IPUrl:      "http://members.3322.org/dyndns/getip",
			LogPath:    "./godns.log",
			LogSize:    16,
			LogNum:     3,
		}

		if err := LoadDomains(os.Getenv("DOMAINS"), &configuration.Domains); err != nil {
			fmt.Println(err.Error())
			log.Println(err.Error())
			os.Exit(1)
		}
	} else {
		//Load settings from configurations file
		if err := LoadSettings(*optConf, &configuration); err != nil {
			fmt.Println(err.Error())
			log.Println(err.Error())
			os.Exit(1)
		}
	}

	if err := checkSettings(&configuration); err != nil {
		log.Println("Settings is invalid! ", err.Error())
		os.Exit(1)
	}

	if err := InitLogger(configuration.LogPath, configuration.LogSize, configuration.LogNum); err != nil {
		log.Println("InitLogger error:", err.Error())
		os.Exit(1)
	}

	dnsLoop()
}

func dnsLoop() {

	for _, domain := range configuration.Domains {
		go domainLoop(&domain)
	}

	select {}
}

func domainLoop(domain *Domain) {
	defer func() {
		if err := recover(); err != nil {
			panicCount++
			log.Printf("Recovered in %v: %v\n", err, debug.Stack())
			fmt.Println(identifyPanic())
			log.Print(identifyPanic())
			if panicCount < PANIC_MAX {
				log.Println("Got panic in goroutine, will start a new one... :", panicCount)
				go domainLoop(domain)
			} else {
				os.Exit(1)
			}
		}
	}()

	for {

		domainID := getDomain(domain.DomainName)

		if domainID == -1 {
			continue
		}

		currentIP, err := getCurrentIP(configuration.IPUrl)

		if err != nil {
			log.Println("get_currentIP:", err)
			continue
		}
		log.Println("currentIp is:", currentIP)

		for _, subDomain := range domain.SubDomains {

			subDomainID, ip := getSubDomain(domainID, subDomain)

			if subDomainID == "" || ip == "" {
				log.Printf("domain: %s.%s subDomainID: %s ip: %s\n", subDomain, domain.DomainName, subDomainID, ip)
				continue
			}

			//Continue to check the IP of sub-domain
			if len(ip) > 0 && !strings.Contains(currentIP, ip) {
				log.Printf("%s.%s Start to update record IP...\n", subDomain, domain.DomainName)
				updateIP(domainID, subDomainID, subDomain, currentIP)
			} else {
				log.Printf("%s.%s Current IP is same as domain IP, no need to update...\n", subDomain, domain.DomainName)
			}
		}

		//Interval is 5 minutes
		time.Sleep(time.Minute * INTERVAL)
	}
}
