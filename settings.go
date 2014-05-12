package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const config_path string = "config.json"

type Settings struct {
	Email      string
	Password   string
	Domain     string
	Sub_domain string
	IP_Url     string
}

func LoadSettings() Settings {
	file, err := ioutil.ReadFile(config_path)

	if err != nil {
		fmt.Println("Error occurs while reading config file, please make sure config file exists!")
		os.Exit(1)
	}

	var setting Settings
	json.Unmarshal(file, &setting)

	return setting
}
