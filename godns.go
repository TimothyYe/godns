package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

const (
	PANIC_MAX         = 5
	DnsUpdateInterval = 5 * time.Minute //Minute
	ClientVersion     = "0.1"           //客户端版本
)

var (
	Configuration  *Settings
	latestIp       string // 上次的IP地址
	optConf        = flag.String("c", "./godns.conf", "config file")
	optCommand     = flag.String("s", "", "send signal to a master process: stop, quit, reopen, reload")
	optHelp        = flag.Bool("h", false, "this help")
	panicCount     = 0
	DomainId       int64
	SubDomainArr   = []string{}
	SubDomainIdArr = []string{}
)

type (
	// 版本
	Version struct {
		ApiVersion    string
		ApiDate       time.Time
		ClientVersion string
	}
	Settings struct {
		ApiId      int    `json:"api_id"`
		ApiToken   string `json:"api_token"`
		Domain     string `json:"domain"`
		SubDomains string `json:"sub_domains"`
		IpFetchUrl string `json:"ip_fetch_url"`
	}
)

func LoadSettings(config_path string) *Settings {
	file, err := ioutil.ReadFile(config_path)
	if err != nil {
		log.Println("Error occurs while reading config file, please make sure config file exists!")
		os.Exit(1)
	}

	var setting Settings
	err = json.Unmarshal(file, &setting)
	if err != nil {
		log.Println("Error occurs while unmarshal config file, please make sure config file correct!")
		os.Exit(1)
	}
	return &setting
}

func usage() {
	log.Println("[command] -c=[config file path]")
	flag.PrintDefaults()
}
func main() {
	flag.Parse()
	if *optHelp {
		usage()
		return
	}
	//log.SetFlags(log.Lshortfile | log.Ltime | log.LstdFlags)

	Configuration = LoadSettings(*optConf)

	ver := GetApiVersion()
	log.Println("[ GoDns][ Version] -", " latest :", ver.ApiVersion,
		" release :", ver.ApiDate.Format("2006-01-02"))
	checkDomain()
	dnsUpdateLoop()
}

// 检测域名
func checkDomain() {
	domain := strings.TrimSpace(Configuration.Domain)
	DomainId = get_domain(domain)
	if DomainId == -1 {
		log.Println("[ GoDns][ Error] - domain :", domain, " dont't be resolve by DnsPod.")
		os.Exit(0)
	}

	SubDomainArr = strings.Split(Configuration.SubDomains, ",")
	for i, v := range SubDomainArr {
		v = strings.TrimSpace(v)
		if len(v) != 0 {
			subDomainId, ip := getSubdomain(DomainId, v)
			subDomain := v + "." + domain
			if subDomainId == "" || ip == "" {
				log.Println("[ GoDns][ Wanning] - ", subDomain, " not in list.")
				SubDomainArr = append(SubDomainArr[:i], SubDomainArr[i+1:]...)
			} else {
				log.Println("[ GoDns][ Stat] - ", subDomain, "=>", ip)
				SubDomainIdArr = append(SubDomainIdArr, subDomainId)
			}
		}
	}
}

func dnsUpdateLoop() {
	defer func() {
		if err := recover(); err != nil {
			panicCount++
			log.Printf("Recovered in %v: %v\n", err, debug.Stack())
			if panicCount < PANIC_MAX {
				log.Println("Got panic in goroutine, will start a new one... :", panicCount)
				go dnsUpdateLoop()
			}
		}
	}()

	for {
		localIp, err := getExternalIp(Configuration.IpFetchUrl)
		if err != nil {
			log.Println("[ GoDns][ Error] - fetch ip error:", err.Error())
			continue
		}
		//检测IP是否有变化,如无变化则不提交更新
		if localIp == latestIp {
			log.Println("[ GoDns][ Stat] - ip not need be update!")
		} else {
			latestIp = localIp
			log.Println("[ GoDns][ Stat] - external ip :", localIp)
			for i, subId := range SubDomainIdArr {
				if err = UpdateIpRecord(DomainId, subId, SubDomainArr[i], localIp); err != nil {
					log.Println("[ GoDns][ Update]- subdomain ", SubDomainArr[i], err.Error())
				} else {
					log.Println("[ GoDns][ Update]- subdomain ", SubDomainArr[i], " update success!")
				}
			}
		}

		//Interval is 5 minutes
		time.Sleep(DnsUpdateInterval)
	}

	log.Printf("Loop %d exited...\n", panicCount)
}
