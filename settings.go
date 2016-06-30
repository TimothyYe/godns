package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

//Settings struct
type Settings struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	LoginToken string `json:"login_token"`
	Domain     string `json:"domain"`
	SubDomain  string `json:"sub_domain"`
	IPUrl      string `json:"ip_url"`
	LogPath    string `json:"log_path"`
	LogSize    int    `json:"log_size"`
	LogNum     int    `json:"log_num"`
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
