package godns

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

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
)

// GetCurrentIP gets public IP from internet
func GetCurrentIP(configuration *Settings) (string, error) {
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
	return string(body), nil
}

// CheckSettings check the format of settings
func CheckSettings(config *Settings) error {
	if config.Provider == DNSPOD {
		if (config.Email == "" || config.Password == "") && config.LoginToken == "" {
			return errors.New("email/password or login token cannot be empty")
		}
	} else if config.Provider == HE {
		if config.Password == "" {
			return errors.New("password cannot be empty")
		}
	} else {
		return errors.New("please provide supported DNS provider: DNSPod/HE")
	}

	return nil
}

// SaveCurrentIP saves current IP into a template file
func SaveCurrentIP(currentIP string) {
	ioutil.WriteFile("./.current_ip", []byte(currentIP), os.FileMode(0644))
}

// LoadCurrentIP loads saved IP from template file
func LoadCurrentIP() string {
	content, err := ioutil.ReadFile("./.current_ip")

	if err != nil {
		return ""
	}

	return strings.Replace(string(content), "\n", "", -1)
}

// SendNotify sends mail notify if IP is changed
func SendNotify(configuration *Settings, domain string, currentIP string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", configuration.Notify.SMTPUsername)
	m.SetHeader("To", configuration.Notify.SendTo)
	m.SetHeader("Subject", "GoDNS Notification")
	log.Println("currentIP:", currentIP)
	log.Println("domain:", domain)
	m.SetBody("text/html", fmt.Sprintf(template, currentIP, domain))

	d := gomail.NewPlainDialer(configuration.Notify.SMTPServer, configuration.Notify.SMTPPort, configuration.Notify.SMTPUsername, configuration.Notify.SMTPPassword)

	// Send the email config by sendlist	.
	if err := d.DialAndSend(m); err != nil {
		log.Println("Send email notification with error:", err.Error())
		return err
	}
	return nil
}
