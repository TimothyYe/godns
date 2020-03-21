package godns

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"

	"golang.org/x/net/proxy"
	gomail "gopkg.in/gomail.v2"
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
	// DNSPOD for dnspod.cn
	DNSPOD = "DNSPod"
	// HE for he.net
	HE = "HE"
	// CLOUDFLARE for cloudflare.com
	CLOUDFLARE = "Cloudflare"
	// ALIDNS for AliDNS
	ALIDNS = "AliDNS"
	// GOOGLE for Google Domains
	GOOGLE = "Google"
	// DUCK for Duck DNS
	DUCK = "DuckDNS"
	// IPV4 for IPV4 mode
	IPV4 = "IPV4"
	// IPV6 for IPV6 mode
	IPV6 = "IPV6"
)

//GetIPFromInterface gets IP address from the specific interface
func GetIPFromInterface(configuration *Settings) (string, error) {
	ifaces, err := net.InterfaceByName(configuration.IPInterface)
	if err != nil {
		log.Println("can't get network device "+configuration.IPInterface+":", err)
		return "", err
	}

	addrs, err := ifaces.Addrs()
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
		if ip == nil {
			continue
		}

		if !(ip.IsGlobalUnicast() &&
			!(ip.IsUnspecified() ||
				ip.IsMulticast() ||
				ip.IsLoopback() ||
				ip.IsLinkLocalUnicast() ||
				ip.IsLinkLocalMulticast() ||
				ip.IsInterfaceLocalMulticast())) {
			continue
		}

		if isIPv4(ip.String()) {
			if strings.ToUpper(configuration.IPType) == IPV4 {
				continue
			}
		} else {
			if strings.ToUpper(configuration.IPType) != IPV6 {
				continue
			}
		}

		return ip.String(), nil

	}
	return "", errors.New("can't get a vaild address from " + configuration.IPInterface)
}

func isIPv4(ip string) bool {
	return strings.Count(ip, ":") < 2
}

// GetHttpClient creates the HTTP client and return it
func GetHttpClient(configuration *Settings) *http.Client {
	client := &http.Client{}

	if configuration.Socks5Proxy != "" {
		log.Println("use socks5 proxy:" + configuration.Socks5Proxy)
		dialer, err := proxy.SOCKS5("tcp", configuration.Socks5Proxy, nil, proxy.Direct)
		if err != nil {
			log.Println("can't connect to the proxy:", err)
			return nil
		}

		httpTransport := &http.Transport{}
		client.Transport = httpTransport
		httpTransport.Dial = dialer.Dial
	}

	return client
}

//GetCurrentIP gets an IP from either internet or specific interface, depending on configuration
func GetCurrentIP(configuration *Settings) (string, error) {
	var err error

	if configuration.IPUrl != "" || configuration.IPV6Url != "" {
		ip, err := GetIPOnline(configuration)
		if err != nil {
			log.Println("get ip online failed. Fallback to get ip from interface if possible.")
		} else {
			return ip, nil
		}
	}

	if configuration.IPInterface != "" {
		ip, err := GetIPFromInterface(configuration)
		if err != nil {
			log.Println("get ip from interface failed. There is no more ways to try.")
		} else {
			return ip, nil
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

	var response *http.Response
	var err error

	if configuration.IPType == "" || strings.ToUpper(configuration.IPType) == IPV4 {
		response, err = client.Get(configuration.IPUrl)
	} else {
		response, err = client.Get(configuration.IPV6Url)
	}

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
		if config.LoginToken == "" {
			if config.Email == "" {
				return errors.New("email cannot be empty")
			}
			if config.Password == "" {
				return errors.New("password cannot be empty")
			}
		}
	} else if config.Provider == ALIDNS {
		if config.Email == "" {
			return errors.New("email cannot be empty")
		}
		if config.Password == "" {
			return errors.New("password cannot be empty")
		}
	} else if config.Provider == DUCK {
		if config.LoginToken == "" {
			return errors.New("login token cannot be empty")
		}
	} else if config.Provider == GOOGLE {
		if config.Email == "" {
			return errors.New("email cannot be empty")
		}
		if config.Password == "" {
			return errors.New("password cannot be empty")
		}
	} else {
		return errors.New("please provide supported DNS provider: DNSPod/HE/AliDNS/Cloudflare/GoogleDomain/DuckDNS")
	}

	return nil
}

// SendNotify sends notify if IP is changed
func SendTelegramNotify(configuration *Settings, domain, currentIP string) error {
	if ! configuration.Notify.Telegram.Enabled {
		return nil
	}

	if configuration.Notify.Telegram.BotApiKey == "" {
		return errors.New("bot api key cannot be empty")
	}

	if configuration.Notify.Telegram.ChatId == "" {
		return errors.New("chat id cannot be empty")
	}


	client := GetHttpClient(configuration)
	tpl := configuration.Notify.Telegram.MsgTemplate
	if tpl == "" {
		tpl = "_Your IP address is changed to_%0A%0A*{{ .CurrentIP }}*%0A%0ADomain *{{ .Domain }}* is updated"
	}

	msg := buildTemplate(currentIP, domain, tpl)
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&parse_mode=Markdown&text=%s",
			configuration.Notify.Telegram.BotApiKey,
			configuration.Notify.Telegram.ChatId,
			msg)
	var response *http.Response
	var err error

	response, err = client.Get(url)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	type ResponseParameters struct {
		MigrateToChatID int64 `json:"migrate_to_chat_id"` // optional
		RetryAfter      int   `json:"retry_after"`        // optional
	}
	type APIResponse struct {
		Ok          bool                `json:"ok"`
		Result      json.RawMessage     `json:"result"`
		ErrorCode   int                 `json:"error_code"`
		Description string              `json:"description"`
		Parameters  *ResponseParameters `json:"parameters"`
	}
	var resp APIResponse
	err = json.Unmarshal([]byte(body), &resp)
	if err != nil {
		fmt.Println("error:", err)
		return errors.New("Failed to parse response")
	}
	if ! resp.Ok {
		return errors.New(resp.Description)
	}

	return nil
}

// SendNotify sends mail notify if IP is changed
func SendMailNotify(configuration *Settings, domain, currentIP string) error {
	if ! configuration.Notify.Mail.Enabled {
		return nil
	}
	log.Print("Sending notification to:", configuration.Notify.Mail.SendTo)
	m := gomail.NewMessage()

	m.SetHeader("From", configuration.Notify.Mail.SMTPUsername)
	m.SetHeader("To", configuration.Notify.Mail.SendTo)
	m.SetHeader("Subject", "GoDNS Notification")
	log.Println("currentIP:", currentIP)
	log.Println("domain:", domain)
	m.SetBody("text/html", buildTemplate(currentIP, domain, mailTemplate))

	d := gomail.NewPlainDialer(configuration.Notify.Mail.SMTPServer, configuration.Notify.Mail.SMTPPort, configuration.Notify.Mail.SMTPUsername, configuration.Notify.Mail.SMTPPassword)

	// Send the email config by sendlist	.
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

// SendNotify sends notify if IP is changed
func SendNotify(configuration *Settings, domain, currentIP string) error {
	err := SendTelegramNotify(configuration, domain, currentIP)
	if (err != nil) {
		log.Println("Send telegram notification with error:", err.Error())
	}
	err = SendMailNotify(configuration, domain, currentIP)
	if (err != nil) {
		log.Println("Send email notification with error:", err.Error())
	}
	return nil
}

func buildTemplate(currentIP, domain string, tplsrc string) string {
	t := template.New("notification template")
	if _, err := t.Parse(tplsrc); err != nil {
		log.Println("Failed to parse template")
		return ""
	}

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
