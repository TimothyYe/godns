package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

//Domain struct
type Domain struct {
	DomainName string   `json:"domain_name"`
	SubDomains []string `json:"sub_domains"`
}

//Settings struct
type Settings struct {
	Provider    string   `json:"provider"`
	Email       string   `json:"email"`
	Password    string   `json:"password"`
	LoginToken  string   `json:"login_token"`
	Domains     []Domain `json:"domains"`
	IPUrl       string   `json:"ip_url"`
	LogPath     string   `json:"log_path"`
	LogSize     int      `json:"log_size"`
	LogNum      int      `json:"log_num"`
	Socks5Proxy string   `json:"socks5_proxy"`
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

//LoadDomains -- Load domains from domains string
func LoadDomains(domainsOrginStr string, domains *[]Domain) error {

	domainsMap := make(map[string]*Domain)
	domainsArray := strings.Split(domainsOrginStr, ",")
	for _, host := range domainsArray {
		dotCount := strings.Count(host, ".")
		if dotCount < 2 {
			continue
		}
		len := len(host)
		pos := strings.Index(host, ".")
		subDomain := host[0:pos]
		domainName := host[pos+1 : len]

		if d, exist := domainsMap[domainName]; exist {
			d.SubDomains = append(d.SubDomains, subDomain)
		} else {
			d := new(Domain)
			d.DomainName = domainName
			d.SubDomains = append(d.SubDomains, subDomain)
			domainsMap[domainName] = d
		}
	}

	for _, d := range domainsMap {
		*domains = append(*domains, *d)
	}

	return nil
}
