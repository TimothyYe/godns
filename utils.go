package godns

import (
	"bytes"
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
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
	// CLOUDFLARE for cloudflare.com
	CLOUDFLARE = "Cloudflare"
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
