package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const (
	// PANIC_MAX is the max allowed panic times
	PANIC_MAX = 5
	// INTERVAL is minute
	INTERVAL = 5
	// DNSPOD
	DNSPOD = "DNSPod"
	// HE
	HE = "HE"
)

var (
	configuration Settings
	optConf       = flag.String("c", "./config.json", "Specify a config file")
	optHelp       = flag.Bool("h", false, "Show help")
	panicCount    = 0
)

func main() {
	flag.Parse()
	if *optHelp {
		flag.Usage()
		return
	}

	//Load settings from configurations file
	if err := LoadSettings(*optConf, &configuration); err != nil {
		fmt.Println(err.Error())
		log.Println(err.Error())
		os.Exit(1)
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
	handler := createHandler(configuration.Provider)
	for _, domain := range configuration.Domains {
		go handler.DomainLoop(&domain)
	}

	select {}
}
