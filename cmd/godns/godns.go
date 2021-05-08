package main

import (
	"flag"
	"os"

	log "github.com/sirupsen/logrus"

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

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	flag.Parse()
	if *optHelp {
		color.Cyan(godns.Logo, Version)
		flag.Usage()
		return
	}

	// Load settings from configurations file
	if err := godns.LoadSettings(*optConf, &configuration); err != nil {
		log.Fatal(err)
	}

	if configuration.DebugInfo {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.DebugLevel)
	}

	if err := godns.CheckSettings(&configuration); err != nil {
		log.Fatal("Settings is invalid! ", err.Error())
	}

	// Init log settings
	log.Info("GoDNS started, entering main loop...")
	dnsLoop()
}

func dnsLoop() {
	panicChan := make(chan godns.Domain)

	log.Info("Creating DNS handler with provider:", configuration.Provider)
	h := handler.CreateHandler(configuration.Provider)
	h.SetConfiguration(&configuration)
	for i := range configuration.Domains {
		go h.DomainLoop(&configuration.Domains[i], panicChan)
	}

	panicCount := 0
	for {
		failDomain := <-panicChan
		log.Debug("Got panic in goroutine, will start a new one... :", panicCount)
		go h.DomainLoop(&failDomain, panicChan)

		panicCount++
		if panicCount >= godns.PanicMax {
			os.Exit(1)
		}
	}
}
