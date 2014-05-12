package main

import (
	"fmt"
	"os"
	"time"
)

var Configuration Settings

func main() {
	fmt.Println("Starting...")

	Configuration = LoadSettings()

	loop := make(chan bool)
	go dns_loop(loop)

	ret := <-loop

	if !ret {
		fmt.Println("Dns loop exited...")
		close(loop)

		os.Exit(1)
	}
}

func dns_loop(loop chan bool) {
	fmt.Println("Inside the loop...")
	time.Sleep(time.Second * 2)

	currentIP, _ := get_currentIP(Configuration.IP_Url)
	fmt.Println("Current IP is" + currentIP)

	//Continue to check the IP of sub-domain
	if len(currentIP) > 0 {
		get_domain(Configuration.Domain)
	}

	api_version()
	loop <- false
}
