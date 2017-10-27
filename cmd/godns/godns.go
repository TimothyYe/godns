package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/TimothyYe/godns"
	"github.com/TimothyYe/godns/handler"
)

const ()

var (
	configuration godns.Settings
	optConf       = flag.String("c", "./config.json", "Specify a config file")
	optHelp       = flag.Bool("h", false, "Show help")
)

func main() {
	flag.Parse()
	if *optHelp {
		flag.Usage()
		return
	}

	// Load settings from configurations file
	if err := godns.LoadSettings(*optConf, &configuration); err != nil {
		fmt.Println(err.Error())
		log.Println(err.Error())
		os.Exit(1)
	}

	if err := godns.CheckSettings(&configuration); err != nil {
		log.Println("Settings is invalid! ", err.Error())
		os.Exit(1)
	}

	if err := godns.InitLogger(configuration.LogPath, configuration.LogSize, configuration.LogNum); err != nil {
		log.Println("InitLogger error:", err.Error())
		os.Exit(1)
	}

	dnsLoop()
}

func dnsLoop() {
	panicChan := make(chan godns.Domain)

	handler := handler.CreateHandler(configuration.Provider)
	handler.SetConfiguration(&configuration)
	for _, domain := range configuration.Domains {
		go handler.DomainLoop(&domain, panicChan)
	}

	panicCount := 0
	for {
		select {
		case failDomain := <-panicChan:
			log.Println("Got panic in goroutine, will start a new one... :", panicCount)
			go handler.DomainLoop(&failDomain, panicChan)
		}

		panicCount++
		if panicCount >= godns.PanicMax {
			os.Exit(1)
		}
	}
}
