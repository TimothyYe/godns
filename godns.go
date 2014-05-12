package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Starting...")

	setting := LoadSettings()
	fmt.Println(setting.IP_Url)

	loop := make(chan bool)
	go dns_loop(setting, loop)

	ret := <-loop

	if !ret {
		fmt.Println("Dns loop exited...")
		close(loop)

		os.Exit(1)
	}
}

func dns_loop(setting Settings, loop chan bool) {

}
