package resolver

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/miekg/dns"
)

func TestNew(t *testing.T) {
	servers := []string{"8.8.8.8", "8.8.4.4"}
	expectedServers := []string{"8.8.8.8:53", "8.8.4.4:53"}
	resolver := New(servers)

	if !reflect.DeepEqual(resolver.Servers, expectedServers) {
		t.Error("resolver.Servers: ", resolver.Servers, "should be equal to", expectedServers)
	}
}

func TestLookupHost_ValidServer(t *testing.T) {
	resolver := New([]string{"8.8.8.8", "8.8.4.4"})
	result, err := resolver.LookupHost("google-public-dns-a.google.com", dns.TypeA)
	if err != nil {
		fmt.Println(err.Error())
		t.Error("Should succeed dns lookup")
	}

	if result[0].String() != "8.8.8.8" {
		t.Error("google-public-dns-a.google.com should be resolved to 8.8.8.8")
	}
}

func TestLookupHostIPv6_ValidServer(t *testing.T) {
	resolver := New([]string{"2001:4860:4860::8888", "2001:4860:4860::8844"})
	result, err := resolver.LookupHost("google-public-dns-a.google.com", dns.TypeAAAA)
	if err != nil {
		fmt.Println(err.Error())
		t.Error("Should succeed dns lookup")
	}

	if result[0].String() != "2001:4860:4860::8888" {
		t.Error("result should be: 2001:4860:4860::8888")
	}
}
