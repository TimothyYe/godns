package provider

import (
	"github.com/TimothyYe/godns/internal/settings"
)

type IDNSProvider interface {
	Init(conf *settings.Settings)
	UpdateIP(domainName string, subdomainName string, ip string) error
}
