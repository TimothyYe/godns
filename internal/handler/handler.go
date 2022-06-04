package handler

import (
	"github.com/TimothyYe/godns/internal/handler/alidns"
	"github.com/TimothyYe/godns/internal/handler/cloudflare"
	"github.com/TimothyYe/godns/internal/handler/dnspod"
	"github.com/TimothyYe/godns/internal/handler/dreamhost"
	"github.com/TimothyYe/godns/internal/handler/duck"
	"github.com/TimothyYe/godns/internal/handler/dynv6"
	"github.com/TimothyYe/godns/internal/handler/google"
	"github.com/TimothyYe/godns/internal/handler/he"
	"github.com/TimothyYe/godns/internal/handler/linode"
	"github.com/TimothyYe/godns/internal/handler/noip"
	"github.com/TimothyYe/godns/internal/handler/scaleway"
	"github.com/TimothyYe/godns/internal/settings"
	"github.com/TimothyYe/godns/internal/utils"
)

// IHandler is the interface for all DNS handlers.
type IHandler interface {
	SetConfiguration(*settings.Settings)
	DomainLoop(domain *settings.Domain, panicChan chan<- settings.Domain, runOnce bool)
}

// CreateHandler creates DNS handler by different providers.
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
	case utils.DYNV6:
		handler = IHandler(&dynv6.Handler{})
	case utils.LINODE:
		handler = IHandler(&linode.Handler{})
	}

	return handler
}
