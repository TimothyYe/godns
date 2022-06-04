package main

import (
	"flag"
	"os"

	"github.com/TimothyYe/godns/internal/handler"
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
		log.Fatal("Settings is invalid! ", err.Error())
	}

	// Init log settings
	log.Info("GoDNS started, entering main loop...")
	dnsLoop()
}

func dnsLoop() {
	panicChan := make(chan settings.Domain)

	log.Infof("Creating DNS handler with provider: %s", configuration.Provider)
	h := handler.CreateHandler(configuration.Provider)
	h.SetConfiguration(&configuration)
	for _, domain := range configuration.Domains {
		if configuration.RunOnce {
			h.DomainLoop(&domain, panicChan, configuration.RunOnce)
		} else {
			go h.DomainLoop(&domain, panicChan, configuration.RunOnce)
		}
	}

	if configuration.RunOnce {
		os.Exit(0)
	}

	panicCount := 0
	for {
		failDomain := <-panicChan
		log.Debug("Got panic in goroutine, will start a new one... :", panicCount)
		go h.DomainLoop(&failDomain, panicChan, configuration.RunOnce)

		panicCount++
		if panicCount >= utils.PanicMax {
			os.Exit(1)
		}
	}
}
