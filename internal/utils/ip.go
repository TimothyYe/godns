package utils

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
	"strings"

	"github.com/TimothyYe/godns/internal/settings"
	dnsResolver "github.com/TimothyYe/godns/pkg/resolver"
	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
)

// GetIPFromInterface gets IP address from the specific interface.
func GetIPFromInterface(configuration *settings.Settings) (string, error) {
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

		if ip.IsPrivate() {
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

// GetCurrentIP gets an IP from either internet or specific interface, depending on configuration.
func GetCurrentIP(configuration *settings.Settings) (string, error) {
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

// GetIPOnline gets public IP from internet.
func GetIPOnline(configuration *settings.Settings) (string, error) {
	client := &http.Client{}
	var reqURL string

	if configuration.IPType == "" || strings.ToUpper(configuration.IPType) == IPV4 {
		reqURL = configuration.IPUrl
	} else {
		reqURL = configuration.IPV6Url
	}

	req, _ := http.NewRequest("GET", reqURL, nil)

	if configuration.UserAgent != "" {
		req.Header.Set("User-Agent", configuration.UserAgent)
	}

	response, err := client.Do(req)

	if err != nil {
		log.Error("Cannot get IP:", err)
		return "", err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get online IP:%d", response.StatusCode)
	}

	body, _ := io.ReadAll(response.Body)
	ipReg := regexp.MustCompile(IPPattern)
	onlineIP := ipReg.FindString(string(body))
	if onlineIP == "" {
		return "", errors.New("failed to get online IP")
	}

	return onlineIP, nil
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
