package utils

import (
	"net"
	"strings"

	dnsResolver "github.com/TimothyYe/godns/pkg/resolver"
	"github.com/miekg/dns"
)

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
		dnsAddress, err := net.LookupHost(hostname)
		if err != nil {
			return "<nil>", err
		}

		return dnsAddress[0], nil
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
