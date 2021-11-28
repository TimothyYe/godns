package handler

import (
	"github.com/TimothyYe/godns/handler/alidns"
	"github.com/TimothyYe/godns/handler/cloudflare"
	"github.com/TimothyYe/godns/handler/dnspod"
	"github.com/TimothyYe/godns/handler/dreamhost"
	"github.com/TimothyYe/godns/handler/duck"
	"github.com/TimothyYe/godns/handler/google"
	"github.com/TimothyYe/godns/handler/he"
	"github.com/TimothyYe/godns/handler/noip"
	"github.com/TimothyYe/godns/handler/scaleway"
	"github.com/TimothyYe/godns/internal/settings"
	"github.com/TimothyYe/godns/internal/utils"
)

// IHandler is the interface for all DNS handlers
type IHandler interface {
	SetConfiguration(*settings.Settings)
	DomainLoop(domain *settings.Domain, panicChan chan<- settings.Domain)
}

// CreateHandler creates DNS handler by different providers
func CreateHandler(provider string) IHandler {
	var handler IHandler

	switch provider {
	case utils.CLOUDFLARE:
		handler = IHandler(&cloudflare.Handler{})
	case utils.DNSPOD:
		handler = IHandler(&dnspod.Handler{})
	case utils.DREAMHOST:
		handler = IHandler(&dreamhost.Handler{})
	case utils.HE:
		handler = IHandler(&he.Handler{})
	case utils.ALIDNS:
		handler = IHandler(&alidns.Handler{})
	case utils.GOOGLE:
		handler = IHandler(&google.Handler{})
	case utils.DUCK:
		handler = IHandler(&duck.Handler{})
	case utils.NOIP:
		handler = IHandler(&noip.Handler{})
	case utils.SCALEWAY:
		handler = IHandler(&scaleway.Handler{})
	}

	return handler
}
