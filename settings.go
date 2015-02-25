package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Settings struct {
	Email      string
	Password   string
	Domain     string
	Sub_domain string
	IP_Url     string
	Log_Path   string
	Log_Size   int
	Log_Num    int
	User       int
	Group      int
}

func LoadSettings(config_path string) Settings {
	file, err := ioutil.ReadFile(config_path)
	if err != nil {
		fmt.Println("Error occurs while reading config file, please make sure config file exists!")
		os.Exit(1)
	}

	var setting Settings
	err = json.Unmarshal(file, &setting)
	if err != nil {
		fmt.Println("Error occurs while unmarshal config file, please make sure config file correct!")
		os.Exit(1)
	}

	return setting
}
