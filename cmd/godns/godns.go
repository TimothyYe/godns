package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/TimothyYe/godns/internal/manager"
	"github.com/TimothyYe/godns/internal/settings"
	"github.com/TimothyYe/godns/internal/utils"

	log "github.com/sirupsen/logrus"

	"github.com/fatih/color"
)

const (
	configEnv = "CONFIG"
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

	configPath := *optConf

	// read config path from the environment
	if os.Getenv(configEnv) != "" {
		// overwrite the config path
		configPath = os.Getenv(configEnv)
	}

	// Load settings from configurations file
	if err := settings.LoadSettings(configPath, &configuration); err != nil {
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

	// Create DNS manager
	dnsManager := &manager.DNSManager{}
	if err := dnsManager.SetConfiguration(&configuration).Build(); err != nil {
		log.Fatal(err)
	}

	// Run DNS manager
	log.Info("GoDNS started, starting the DNS manager...")
	dnsManager.Run()

	// handle the signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// stop the DNS manager
	<-c
	log.Info("GoDNS is terminated, stopping the DNS manager...")
	dnsManager.Stop()
	log.Info("GoDNS is stopped, bye!")
}
