package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

//Settings struct
type Settings struct {
	Email      string
	Password   string
	LoginToken string
	Domain     string
	Sub_domain string
	IP_Url     string
	Log_Path   string
	Log_Size   int
	Log_Num    int
}

//LoadSettings -- Load settings from config file
func LoadSettings(configPath string, settings *Settings) error {
	//LoadSettings from config file
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Println("Error occurs while reading config file, please make sure config file exists!")
		return err
	}

	err = json.Unmarshal(file, settings)
	if err != nil {
		fmt.Println("Error occurs while unmarshal config file, please make sure config file correct!")
		return err
	}

	return nil
}
