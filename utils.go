package godns

import (
	"bytes"
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
    "net"

	"golang.org/x/net/proxy"
	"gopkg.in/gomail.v2"
)

var (
	// Logo for GoDNS
	Logo = `

 ██████╗  ██████╗ ██████╗ ███╗   ██╗███████╗
██╔════╝ ██╔═══██╗██╔══██╗████╗  ██║██╔════╝
██║  ███╗██║   ██║██║  ██║██╔██╗ ██║███████╗
██║   ██║██║   ██║██║  ██║██║╚██╗██║╚════██║
╚██████╔╝╚██████╔╝██████╔╝██║ ╚████║███████║
 ╚═════╝  ╚═════╝ ╚═════╝ ╚═╝  ╚═══╝╚══════╝

GoDNS V%s
https://github.com/TimothyYe/godns

`
)

const (
	// PanicMax is the max allowed panic times
	PanicMax = 5
	// INTERVAL is minute
	INTERVAL = 5
	// DNSPOD for dnspod.cn
	DNSPOD = "DNSPod"
	// HE for he.net
	HE = "HE"
	// CLOUDFLARE for cloudflare.com
	CLOUDFLARE = "Cloudflare"
)

//GetIPFromInterface gets IP address from the specific interface
func GetIPFromInterface(configuration *Settings) (string,error) {
 	ifaces, err := net.InterfaceByName(configuration.IPInterface)
	if err != nil {
		log.Println("can't get network device "+configuration.IPInterface+":", err)
		return "", err
	}

	addrs, err  :=ifaces.Addrs();
	if err != nil {
		log.Println("can't get address from "+configuration.IPInterface+":", err)
		return "", err
	}

	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}
		if (ip == nil){
			continue;
		}

		if (!(ip.IsGlobalUnicast() &&!(ip.IsUnspecified()||ip.IsMulticast()||ip.IsLoopback()||ip.IsLinkLocalUnicast()||ip.IsLinkLocalMulticast()||ip.IsInterfaceLocalMulticast()))){
			continue;
		} 

		//the code is not ready for updating an AAAA record
		/*
		if (isIPv4(ip.String())){
			if (configuration.IPType=="IPv6"){
				continue;
			}
		}else{
			if (configuration.IPType!="IPv6"){
				continue;
			}
		} */
		if (!isIPv4(ip.String())){
			continue;
		}

		return ip.String(),nil
		 
	}
	return "", errors.New("can't get a vaild address from "+configuration.IPInterface)
}

func isIPv4(ip string) bool{
	return strings.Count(ip, ":") < 2
}

//GetCurrentIP gets an IP from either internet or specific interface, depending on configuration
func GetCurrentIP(configuration *Settings) (string, error) {
	var err error

	if (configuration.IPUrl != ""){ 
		ip, err := GetIPOnline(configuration)
		if (err != nil){
			log.Println("get ip online failed. Fallback to get ip from interface if possible.")
		}else{
			return ip,nil
		}
	}

	if (configuration.IPInterface != ""){
		ip, err := GetIPFromInterface(configuration)
		if (err != nil){ 
			log.Println("get ip from interface failed. There is no more ways to try.")
		}else{
			return ip,nil
		}
	}

	return "", err
}
// GetIPOnline gets public IP from internet
func GetIPOnline(configuration *Settings) (string, error) {
	client := &http.Client{}

	if configuration.Socks5Proxy != "" {

		log.Println("use socks5 proxy:" + configuration.Socks5Proxy)
		dialer, err := proxy.SOCKS5("tcp", configuration.Socks5Proxy, nil, proxy.Direct)
		if err != nil {
			log.Println("can't connect to the proxy:", err)
			return "", err
		}

		httpTransport := &http.Transport{}
		client.Transport = httpTransport
		httpTransport.Dial = dialer.Dial
	}

	response, err := client.Get(configuration.IPUrl)

	if err != nil {
		log.Println("Cannot get IP...")
		return "", err
	}

	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	return strings.Trim(string(body), "\n"), nil
}

// CheckSettings check the format of settings
func CheckSettings(config *Settings) error {
	if config.Provider == DNSPOD {
		if config.Password == "" && config.LoginToken == "" {
			return errors.New("password or login token cannot be empty")
		}
	} else if config.Provider == HE {
		if config.Password == "" {
			return errors.New("password cannot be empty")
		}
	} else if config.Provider == CLOUDFLARE {
		if config.Email == "" {
			return errors.New("email cannot be empty")
		}
		if config.Password == "" {
			return errors.New("password cannot be empty")
		}
	} else {
		return errors.New("please provide supported DNS provider: DNSPod/HE")
	}

	return nil
}

// SendNotify sends mail notify if IP is changed
func SendNotify(configuration *Settings, domain, currentIP string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", configuration.Notify.SMTPUsername)
	m.SetHeader("To", configuration.Notify.SendTo)
	m.SetHeader("Subject", "GoDNS Notification")
	log.Println("currentIP:", currentIP)
	log.Println("domain:", domain)
	m.SetBody("text/html", buildTemplate(currentIP, domain))

	d := gomail.NewPlainDialer(configuration.Notify.SMTPServer, configuration.Notify.SMTPPort, configuration.Notify.SMTPUsername, configuration.Notify.SMTPPassword)

	// Send the email config by sendlist	.
	if err := d.DialAndSend(m); err != nil {
		log.Println("Send email notification with error:", err.Error())
		return err
	}
	return nil
}

func buildTemplate(currentIP, domain string) string {
	t := template.New("notification template")
	t.Parse(mailTemplate)

	data := struct {
		CurrentIP string
		Domain    string
	}{
		currentIP,
		domain,
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		log.Println(err.Error())
		return ""
	}

	return tpl.String()
}
