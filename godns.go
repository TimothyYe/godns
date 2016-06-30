package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
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

func identifyPanic() string {
	var name, file string
	var line int
	var pc [16]uintptr

	n := runtime.Callers(3, pc[:])
	for _, pc := range pc[:n] {
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}
		file, line = fn.FileLine(pc)
		name = fn.Name()
		if !strings.HasPrefix(name, "runtime.") {
			break
		}
	}

	switch {
	case name != "":
		return fmt.Sprintf("%v:%v", name, line)
	case file != "":
		return fmt.Sprintf("%v:%v", file, line)
	}

	return fmt.Sprintf("pc:%x", pc)
}

func usage() {
	log.Println("[command] -c=[config file path]")
	flag.PrintDefaults()
}

func main() {
	flag.Parse()
	if *optHelp {
		usage()
		return
	}

	if err := LoadSettings(*optConf, &configuration); err != nil {
		fmt.Println(err.Error())
		log.Println(err.Error())
		os.Exit(1)
	}

	if err := InitLogger(configuration.Log_Path, configuration.Log_Size, configuration.Log_Num); err != nil {
		log.Println("InitLogger error:", err)
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

		currentIP, err := getCurrentIP(configuration.IP_Url)

		if err != nil {
			log.Println("get_currentIP:", err)
			continue
		}

		subDomainID, ip := getSubDomain(domainID, configuration.Sub_domain)

		if subDomainID == "" || ip == "" {
			log.Println("sub_domain:", subDomainID, ip)
			continue
		}

		log.Println("currentIp is:", currentIP)

		//Continue to check the IP of sub-domain
		if len(ip) > 0 && !strings.Contains(currentIP, ip) {
			log.Println("Start to update record IP...")
			updateIP(domainID, subDomainID, configuration.Sub_domain, currentIP)
		} else {
			log.Println("Current IP is same as domain IP, no need to update...")
		}

		//Interval is 5 minutes
		time.Sleep(time.Minute * INTERVAL)
	}
}
