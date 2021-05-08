package godns

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"golang.org/x/net/proxy"

	dnsResolver "github.com/TimothyYe/godns/resolver"

	"github.com/miekg/dns"
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
	// DREAMHOST for Dreamhost
	DREAMHOST = "Dreamhost"
	// NOIP for NoIP
	NOIP = "NoIP"
	// IPV4 for IPV4 mode
	IPV4 = "IPV4"
	// IPV6 for IPV6 mode
	IPV6 = "IPV6"
	// IPTypeA
	IPTypeA = "A"
	// IPTypeAAAA
	IPTypeAAAA = "AAAA"
	// RootDomain
	RootDomain = "@"
)

//GetIPFromInterface gets IP address from the specific interface
func GetIPFromInterface(configuration *Settings) (string, error) {
	ifaces, err := net.InterfaceByName(configuration.IPInterface)
	if err != nil {
		log.Error("can't get network device "+configuration.IPInterface+":", err)
		return "", err
	}

	addrs, err := ifaces.Addrs()
	if err != nil {
		log.Error("can't get address from "+configuration.IPInterface+":", err)
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
			if strings.ToUpper(configuration.IPType) != IPV4 {
				continue
			}
		} else {
			if strings.ToUpper(configuration.IPType) != IPV6 {
				continue
			}
		}

		if ip.String() != "" {
			return ip.String(), nil
		}
	}
	return "", errors.New("can't get a vaild address from " + configuration.IPInterface)
}

func isIPv4(ip string) bool {
	return strings.Count(ip, ":") < 2
}

//GetCurrentIP gets an IP from either internet or specific interface, depending on configuration
func GetCurrentIP(configuration *Settings) (string, error) {
	var err error
	var ip string

	if configuration.IPUrl != "" || configuration.IPV6Url != "" {
		ip, err = GetIPOnline(configuration)
		if err != nil {
			log.Error("get ip online failed. Fallback to get ip from interface if possible.")
		} else {
			return ip, nil
		}
	}

	if configuration.IPInterface != "" {
		ip, err = GetIPFromInterface(configuration)
		if err != nil {
			log.Error("get ip from interface failed. There is no more ways to try.")
		} else {
			return ip, nil
		}
	}

	return "", err
}

// GetIPOnline gets public IP from internet
func GetIPOnline(configuration *Settings) (string, error) {
	client := &http.Client{}

	var response *http.Response
	var err error

	if configuration.IPType == "" || strings.ToUpper(configuration.IPType) == IPV4 {
		response, err = client.Get(configuration.IPUrl)
	} else {
		response, err = client.Get(configuration.IPV6Url)
	}

	if err != nil {
		log.Error("Cannot get IP:", err)
		return "", err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get online IP:%d", response.StatusCode)
	}

	body, _ := ioutil.ReadAll(response.Body)
	onlineIP := strings.Trim(string(body), "\n")
	if onlineIP == "" {
		return "", errors.New("failed to get online IP")
	}

	return onlineIP, nil
}

// CheckSettings check the format of settings
func CheckSettings(config *Settings) error {
	switch config.Provider {
	case DNSPOD:
		if config.Password == "" && config.LoginToken == "" {
			return errors.New("password or login token cannot be empty")
		}
	case HE:
		if config.Password == "" {
			return errors.New("password cannot be empty")
		}
	case CLOUDFLARE:
		if config.LoginToken == "" {
			if config.Email == "" {
				return errors.New("email cannot be empty")
			}
			if config.Password == "" {
				return errors.New("password cannot be empty")
			}
		}
	case ALIDNS:
		if config.Email == "" {
			return errors.New("email cannot be empty")
		}
		if config.Password == "" {
			return errors.New("password cannot be empty")
		}
	case DUCK:
		if config.LoginToken == "" {
			return errors.New("login token cannot be empty")
		}
	case GOOGLE:
		fallthrough
	case NOIP:
		if config.Email == "" {
			return errors.New("email cannot be empty")
		}
		if config.Password == "" {
			return errors.New("password cannot be empty")
		}
	case DREAMHOST:
		if config.LoginToken == "" {
			return errors.New("login token cannot be empty")
		}

	default:
		return errors.New("please provide supported DNS provider: DNSPod/HE/AliDNS/Cloudflare/GoogleDomain/DuckDNS/Dreamhost")

	}

	return nil
}

// GetHttpClient creates the HTTP client and return it
func GetHttpClient(conf *Settings, useProxy bool) *http.Client {
	client := &http.Client{}

	if useProxy && conf.Socks5Proxy != "" {
		log.Debug("use socks5 proxy:" + conf.Socks5Proxy)
		dialer, err := proxy.SOCKS5("tcp", conf.Socks5Proxy, nil, proxy.Direct)
		if err != nil {
			log.Error("can't connect to the proxy:", err)
			return nil
		}

		httpTransport := &http.Transport{}
		client.Transport = httpTransport
		httpTransport.Dial = dialer.Dial
	}

	return client
}

// ResolveDNS will query DNS for a given hostname.
func ResolveDNS(hostname, resolver, ipType string) (string, error) {
	var dnsType uint16
	if ipType == "" || strings.ToUpper(ipType) == IPV4 {
		dnsType = dns.TypeA
	} else {
		dnsType = dns.TypeAAAA
	}

	// If no DNS server is set in config file, falls back to default resolver.
	if resolver == "" {
		dnsAdress, err := net.LookupHost(hostname)
		if err != nil {
			return "<nil>", err
		}

		return dnsAdress[0], nil
	}
	res := dnsResolver.New([]string{resolver})
	// In case of i/o timeout
	res.RetryTimes = 5

	ip, err := res.LookupHost(hostname, dnsType)
	if err != nil {
		return "<nil>", err
	}

	return ip[0].String(), nil

}
