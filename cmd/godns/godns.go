package main

import (
	"flag"
	"fmt"
	"os"

	"log"

	"github.com/TimothyYe/godns"
	"github.com/TimothyYe/godns/handler"
	"github.com/fatih/color"
)

var (
	configuration godns.Settings
	optConf       = flag.String("c", "./config.json", "Specify a config file")
	optHelp       = flag.Bool("h", false, "Show help")

	// Version is current version of GoDNS
	Version = "0.1"
)

func main() {
	flag.Parse()
	if *optHelp {
		color.Cyan(godns.Logo, Version)
		flag.Usage()
		return
	}

	// Load settings from configurations file
	if err := godns.LoadSettings(*optConf, &configuration); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if err := godns.CheckSettings(&configuration); err != nil {
		fmt.Println("Settings is invalid! ", err.Error())
		os.Exit(1)
	}

	// Init log settings
	log.SetPrefix("【GoDNS】")
	log.Println("GoDNS started, entering main loop...")
	dnsLoop()
}

func dnsLoop() {
	panicChan := make(chan godns.Domain)

	log.Println("Creating DNS handler with provider:", configuration.Provider)
	handler := handler.CreateHandler(configuration.Provider)
	handler.SetConfiguration(&configuration)
	for i, _ := range configuration.Domains {
		go handler.DomainLoop(&configuration.Domains[i], panicChan)
	}

	panicCount := 0
	for {
		failDomain := <-panicChan
		log.Println("Got panic in goroutine, will start a new one... :", panicCount)
		go handler.DomainLoop(&failDomain, panicChan)

		panicCount++
		if panicCount >= godns.PanicMax {
			os.Exit(1)
		}
	}
}
