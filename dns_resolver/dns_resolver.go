// Package dns_resolver is a simple dns resolver
// based on miekg/dns
package dns_resolver

import (
	"errors"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"

	"github.com/miekg/dns"
)

// DnsResolver represents a dns resolver
type DnsResolver struct {
	Servers    []string
	RetryTimes int
	r          *rand.Rand
}

// New initializes DnsResolver.
func New(servers []string) *DnsResolver {
	for i := range servers {
		servers[i] = net.JoinHostPort(servers[i], "53")
	}

	return &DnsResolver{servers, len(servers) * 2, rand.New(rand.NewSource(time.Now().UnixNano()))}
}

// NewFromResolvConf initializes DnsResolver from resolv.conf like file.
func NewFromResolvConf(path string) (*DnsResolver, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &DnsResolver{}, errors.New("no such file or directory: " + path)
	}
	config, err := dns.ClientConfigFromFile(path)
	servers := []string{}
	for _, ipAddress := range config.Servers {
		servers = append(servers, net.JoinHostPort(ipAddress, "53"))
	}
	return &DnsResolver{servers, len(servers) * 2, rand.New(rand.NewSource(time.Now().UnixNano()))}, err
}

// LookupHost returns IP addresses of provied host.
// In case of timeout retries query RetryTimes times.
func (r *DnsResolver) LookupHost(host string, dnsType uint16) ([]net.IP, error) {
	return r.lookupHost(host, dnsType, r.RetryTimes)
}

func (r *DnsResolver) lookupHost(host string, dnsType uint16, triesLeft int) ([]net.IP, error) {
	m1 := new(dns.Msg)
	m1.Id = dns.Id()
	m1.RecursionDesired = true
	m1.Question = make([]dns.Question, 1)

	switch dnsType {
	case dns.TypeA:
		m1.Question[0] = dns.Question{dns.Fqdn(host), dns.TypeA, dns.ClassINET}
	case dns.TypeAAAA:
		m1.Question[0] = dns.Question{dns.Fqdn(host), dns.TypeAAAA, dns.ClassINET}
	}

	in, err := dns.Exchange(m1, r.Servers[r.r.Intn(len(r.Servers))])

	result := []net.IP{}

	if err != nil {
		if strings.HasSuffix(err.Error(), "i/o timeout") && triesLeft > 0 {
			triesLeft--
			return r.lookupHost(host, dnsType, triesLeft)
		}
		return result, err
	}

	if in != nil && in.Rcode != dns.RcodeSuccess {
		return result, errors.New(dns.RcodeToString[in.Rcode])
	}

	if dnsType == dns.TypeA {
		for _, record := range in.Answer {
			if t, ok := record.(*dns.A); ok {
				result = append(result, t.A)
			}
		}
	}

	if dnsType == dns.TypeAAAA {
		for _, record := range in.Answer {
			if t, ok := record.(*dns.AAAA); ok {
				result = append(result, t.AAAA)
			}
		}
	}

	return result, err
}
