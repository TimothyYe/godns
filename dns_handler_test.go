package main

import (
	"testing"
)

func Test_get_current_IP(t *testing.T) {
	ip, _ := get_currentIP("http://members.3322.org/dyndns/getip")

	if ip == "" {
		t.Log("IP is empty...")
	} else {
		t.Log("IP is:" + ip)
	}
}
