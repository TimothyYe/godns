package handler

import (
	"github.com/TimothyYe/godns/internal/handler/cloudflare"
	"github.com/TimothyYe/godns/internal/handler/duck"
	"github.com/TimothyYe/godns/internal/handler/google"
	"github.com/TimothyYe/godns/internal/handler/he"
	"github.com/TimothyYe/godns/internal/handler/noip"
	"github.com/TimothyYe/godns/internal/provider/alidns"
	"github.com/TimothyYe/godns/internal/provider/dnspod"
	"github.com/TimothyYe/godns/internal/provider/dreamhost"
	"github.com/TimothyYe/godns/internal/provider/dynv6"
	"github.com/TimothyYe/godns/internal/provider/linode"
	"github.com/TimothyYe/godns/internal/provider/scaleway"
	"github.com/TimothyYe/godns/internal/settings"
	"github.com/TimothyYe/godns/internal/utils"
)

// IHandler is the interface for all DNS handlers.
type IHandler interface {
	SetConfiguration(*settings.Settings)
	DomainLoop(domain *settings.Domain, panicChan chan<- settings.Domain, runOnce bool)
}

// CreateHandler creates DNS handler by different providers.
func CreateHandler(conf *settings.Settings) (IHandler, error) {
	var handler IHandler

	switch conf.Provider {
	case utils.CLOUDFLARE:
		handler = IHandler(&cloudflare.Handler{})
	case utils.DNSPOD:
		dnsPodProvider := dnspod.DNSProvider{}
		dnsPodProvider.Init(conf)
		genericHandler := Handler{}
		genericHandler.SetProvider(&dnsPodProvider)
		handler = IHandler(&genericHandler)
	case utils.DREAMHOST:
		dreamHostProvider := dreamhost.DNSProvider{}
		dreamHostProvider.Init(conf)
		genericHandler := Handler{}
		genericHandler.SetProvider(&dreamHostProvider)
		handler = IHandler(&genericHandler)
	case utils.HE:
		handler = IHandler(&he.Handler{})
	case utils.ALIDNS:
		aliDNSProvider := alidns.DNSProvider{}
		aliDNSProvider.Init(conf)
		genericHandler := Handler{}
		genericHandler.SetProvider(&aliDNSProvider)
		handler = IHandler(&genericHandler)
	case utils.GOOGLE:
		handler = IHandler(&google.Handler{})
	case utils.DUCK:
		handler = IHandler(&duck.Handler{})
	case utils.NOIP:
		handler = IHandler(&noip.Handler{})
	case utils.SCALEWAY:
		scaleWayProvider := scaleway.DNSProvider{}
		scaleWayProvider.Init(conf)
		genericHandler := Handler{}
		genericHandler.SetProvider(&scaleWayProvider)
		handler = IHandler(&genericHandler)
	case utils.DYNV6:
		dynV6Provider := dynv6.DNSProvider{}
		dynV6Provider.Init(conf)
		genericHandler := Handler{}
		genericHandler.SetProvider(&dynV6Provider)
		handler = IHandler(&genericHandler)
	case utils.LINODE:
		linodeProvider := linode.DNSProvider{}
		linodeProvider.Init(conf)
		genericHandler := Handler{}
		genericHandler.SetProvider(&linodeProvider)
		handler = IHandler(&genericHandler)
	default:
		return nil, utils.ErrUnknownProvider
	}

	return handler, nil
}
