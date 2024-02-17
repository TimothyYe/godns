package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	config  settings.Settings
	optAddr = flag.String("a", ":9000", "Specify the address to listen on")
	optConf = flag.String("c", "./config.json", "Specify a config file")
	optHelp = flag.Bool("h", false, "Show help")

	// Version is current version of GoDNS.
	Version = "v0.1"
)

func main() {
	utils.Version = Version

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

	// Load settings from configs file
	if err := settings.LoadSettings(configPath, &config); err != nil {
		log.Fatal(err)
	}

	// set the log level
	log.SetOutput(os.Stdout)

	if config.DebugInfo {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if err := utils.CheckSettings(&config); err != nil {
		log.Fatal("Invalid settings: ", err.Error())
	}

	// Create DNS manager
	dnsManager := manager.GetDNSManager(configPath, &config, *optAddr)

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
	// wait for the goroutines to exit
	time.Sleep(200 * time.Millisecond)
	log.Info("GoDNS is stopped, bye!")
}
