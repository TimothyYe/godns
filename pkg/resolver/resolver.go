// Package resolver is a simple dns resolver
// based on miekg/dns
package resolver

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"

	"github.com/miekg/dns"
)

// DNSResolver represents a dns resolver.
type DNSResolver struct {
	Servers    []string
	RetryTimes int
	r          *rand.Rand
}

// New initializes DnsResolver.
func New(servers []string) *DNSResolver {
	for i := range servers {
		servers[i] = net.JoinHostPort(servers[i], "53")
	}

	return &DNSResolver{servers, len(servers) * 2, rand.New(rand.NewSource(time.Now().UnixNano()))}
}

// NewFromResolvConf initializes DnsResolver from resolv.conf like file.
func NewFromResolvConf(path string) (*DNSResolver, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &DNSResolver{}, errors.New("no such file or directory: " + path)
	}

	config, err := dns.ClientConfigFromFile(path)
	if err != nil {
		return &DNSResolver{}, err
	}

	var servers []string
	for _, ipAddress := range config.Servers {
		servers = append(servers, net.JoinHostPort(ipAddress, "53"))
	}
	return &DNSResolver{servers, len(servers) * 2, rand.New(rand.NewSource(time.Now().UnixNano()))}, err
}

// LookupHost returns IP addresses of provied host.
// In case of timeout retries query RetryTimes times.
func (r *DNSResolver) LookupHost(host string, dnsType uint16) ([]net.IP, error) {
	return r.lookupHost(host, dnsType, r.RetryTimes)
}

func (r *DNSResolver) lookupHost(host string, dnsType uint16, triesLeft int) ([]net.IP, error) {
	m1 := new(dns.Msg)
	m1.Id = dns.Id()
	m1.RecursionDesired = true
	m1.Question = make([]dns.Question, 1)

	switch dnsType {
	case dns.TypeA:
		m1.Question[0] = dns.Question{Name: dns.Fqdn(host), Qtype: dns.TypeA, Qclass: dns.ClassINET}
	case dns.TypeAAAA:
		m1.Question[0] = dns.Question{Name: dns.Fqdn(host), Qtype: dns.TypeAAAA, Qclass: dns.ClassINET}
	}

	in, err := dns.Exchange(m1, r.Servers[r.r.Intn(len(r.Servers))])

	var result []net.IP

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
		if len(in.Answer) > 0 {
			for _, record := range in.Answer {
				if t, ok := record.(*dns.A); ok {
					result = append(result, t.A)
				}
			}
		} else {
			return result, errors.New("empty result")
		}
	}

	if dnsType == dns.TypeAAAA {
		if len(in.Answer) > 0 {
			for _, record := range in.Answer {
				if t, ok := record.(*dns.AAAA); ok {
					result = append(result, t.AAAA)
				}
			}
		} else {
			return result, fmt.Errorf("cannot resolve domain %s, please make sure the IP type is right", host)
		}
	}

	return result, err
}
