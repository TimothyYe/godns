package main

import (
	"flag"
	"os"

	"github.com/TimothyYe/godns/internal/handler"
	"github.com/TimothyYe/godns/internal/provider"
	"github.com/TimothyYe/godns/internal/settings"
	"github.com/TimothyYe/godns/internal/utils"

	log "github.com/sirupsen/logrus"

	"github.com/fatih/color"
)

var (
	configuration settings.Settings
	optConf       = flag.String("c", "./config.json", "Specify a config file")
	optHelp       = flag.Bool("h", false, "Show help")

	// Version is current version of GoDNS.
	Version = "0.1"
)

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	flag.Parse()
	if *optHelp {
		color.Cyan(utils.Logo, Version)
		flag.Usage()
		return
	}

	// Load settings from configurations file
	if err := settings.LoadSettings(*optConf, &configuration); err != nil {
		log.Fatal(err)
	}

	if configuration.DebugInfo {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if err := utils.CheckSettings(&configuration); err != nil {
		log.Fatal("Invalid settings: ", err.Error())
	}

	// Init log settings
	log.Info("GoDNS started, entering main loop...")
	dnsLoop()
}

func dnsLoop() {
	panicChan := make(chan settings.Domain)

	log.Infof("Creating DNS handler with provider: %s", configuration.Provider)
	dnsProvider, err := provider.GetProvider(&configuration)
	if err != nil {
		log.Fatal(err)
	}

	ddnsHandler := handler.Handler{}
	ddnsHandler.SetConfiguration(&configuration)
	ddnsHandler.SetProvider(dnsProvider)

	for _, domain := range configuration.Domains {
		domain := domain
		if configuration.RunOnce {
			err := ddnsHandler.UpdateIP(&domain)
			if err != nil {
				log.Error("Error during execution:", err)
				os.Exit(1)
			}
		} else {
			go ddnsHandler.LoopUpdateIP(&domain, panicChan)
		}
	}

	if configuration.RunOnce {
		os.Exit(0)
	}

	panicCount := 0
	for {
		failDomain := <-panicChan
		log.Debug("Got panic in goroutine, will start a new one... :", panicCount)
		go ddnsHandler.LoopUpdateIP(&failDomain, panicChan)

		panicCount++
		if panicCount >= utils.PanicMax {
			os.Exit(1)
		}
	}
}
