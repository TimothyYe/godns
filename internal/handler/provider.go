package handler

import (
	"github.com/TimothyYe/godns/internal/settings"
)

type DNSProvider interface {
	Init(conf *settings.Settings) error
	UpdateIP(domainName string, subdomainName string, ip string) error
}
