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
	PANIC_MAX = 5
	INTERVAL  = 5 //Minute
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
			Domain:     os.Getenv("DOMAIN"),
			SubDomain:  os.Getenv("SUB_DOMAIN"),
			IPUrl:      "http://members.3322.org/dyndns/getip",
			LogPath:    "./godns.log",
			LogSize:    16,
			LogNum:     3,
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
		log.Println("Settings is invalid! " + err.Error())
		os.Exit(1)
	}

	if err := InitLogger(configuration.LogPath, configuration.LogSize, configuration.LogNum); err != nil {
		log.Println("InitLogger error:" + err.Error())
		os.Exit(1)
	}

	dnsLoop()
}

func dnsLoop() {
	defer func() {
		if err := recover(); err != nil {
			panicCount++
			log.Printf("Recovered in %v: %v\n", err, debug.Stack())
			fmt.Println(identifyPanic())
			log.Print(identifyPanic())
			if panicCount < PANIC_MAX {
				log.Println("Got panic in goroutine, will start a new one... :", panicCount)
				go dnsLoop()
			}
		}
	}()

	for {

		domainID := getDomain(configuration.Domain)

		if domainID == -1 {
			continue
		}

		currentIP, err := getCurrentIP(configuration.IPUrl)

		if err != nil {
			log.Println("get_currentIP:", err)
			continue
		}

		subDomainID, ip := getSubDomain(domainID, configuration.SubDomain)

		if subDomainID == "" || ip == "" {
			log.Println("sub_domain:", subDomainID, ip)
			continue
		}

		log.Println("currentIp is:", currentIP)

		//Continue to check the IP of sub-domain
		if len(ip) > 0 && !strings.Contains(currentIP, ip) {
			log.Println("Start to update record IP...")
			updateIP(domainID, subDomainID, configuration.SubDomain, currentIP)
		} else {
			log.Println("Current IP is same as domain IP, no need to update...")
		}

		//Interval is 5 minutes
		time.Sleep(time.Minute * INTERVAL)
	}
}
