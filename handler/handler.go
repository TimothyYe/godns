package handler

import (
	"github.com/TimothyYe/godns"
	"github.com/TimothyYe/godns/handler/alidns"
	"github.com/TimothyYe/godns/handler/cloudflare"
	"github.com/TimothyYe/godns/handler/dnspod"
	"github.com/TimothyYe/godns/handler/he"
)

// IHandler is the interface for all DNS handlers
type IHandler interface {
	SetConfiguration(*godns.Settings)
	DomainLoop(domain *godns.Domain, panicChan chan<- godns.Domain)
}

// CreateHandler creates DNS handler by different providers
func CreateHandler(provider string) IHandler {
	var handler IHandler

	switch provider {
	case godns.CLOUDFLARE:
		handler = IHandler(&cloudflare.CloudflareHandler{})
	case godns.DNSPOD:
		handler = IHandler(&dnspod.DNSPodHandler{})
	case godns.HE:
		handler = IHandler(&he.HEHandler{})
	case godns.ALIDNS:
		handler = IHandler(&alidns.AliDNSHandler{})
	}

	return handler
}
