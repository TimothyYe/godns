package handler

import (
	"github.com/TimothyYe/godns/internal/provider/alidns"
	"github.com/TimothyYe/godns/internal/provider/cloudflare"
	"github.com/TimothyYe/godns/internal/provider/dnspod"
	"github.com/TimothyYe/godns/internal/provider/dreamhost"
	"github.com/TimothyYe/godns/internal/provider/duck"
	"github.com/TimothyYe/godns/internal/provider/dynv6"
	"github.com/TimothyYe/godns/internal/provider/google"
	"github.com/TimothyYe/godns/internal/provider/he"
	"github.com/TimothyYe/godns/internal/provider/linode"
	"github.com/TimothyYe/godns/internal/provider/noip"
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
	genericHandler := Handler{}

	switch conf.Provider {
	case utils.CLOUDFLARE:
		cfDNSProvider := cloudflare.DNSProvider{}
		cfDNSProvider.Init(conf)
		genericHandler.SetProvider(&cfDNSProvider)
	case utils.DNSPOD:
		dnsPodProvider := dnspod.DNSProvider{}
		dnsPodProvider.Init(conf)
		genericHandler.SetProvider(&dnsPodProvider)
	case utils.DREAMHOST:
		dreamHostProvider := dreamhost.DNSProvider{}
		dreamHostProvider.Init(conf)
		genericHandler.SetProvider(&dreamHostProvider)
	case utils.HE:
		heProvider := he.DNSProvider{}
		heProvider.Init(conf)
		genericHandler.SetProvider(&heProvider)
	case utils.ALIDNS:
		aliDNSProvider := alidns.DNSProvider{}
		aliDNSProvider.Init(conf)
		genericHandler.SetProvider(&aliDNSProvider)
	case utils.GOOGLE:
		googleDNSProvider := google.DNSProvider{}
		googleDNSProvider.Init(conf)
		genericHandler.SetProvider(&googleDNSProvider)
	case utils.DUCK:
		duckDNSProvider := duck.DNSProvider{}
		duckDNSProvider.Init(conf)
		genericHandler.SetProvider(&duckDNSProvider)
	case utils.NOIP:
		noIPProvider := noip.DNSProvider{}
		noIPProvider.Init(conf)
		genericHandler.SetProvider(&noIPProvider)
	case utils.SCALEWAY:
		scaleWayProvider := scaleway.DNSProvider{}
		scaleWayProvider.Init(conf)
		genericHandler.SetProvider(&scaleWayProvider)
	case utils.DYNV6:
		dynV6Provider := dynv6.DNSProvider{}
		dynV6Provider.Init(conf)
		genericHandler.SetProvider(&dynV6Provider)
	case utils.LINODE:
		linodeProvider := linode.DNSProvider{}
		linodeProvider.Init(conf)
		genericHandler.SetProvider(&linodeProvider)
	default:
		return nil, utils.ErrUnknownProvider
	}

	handler = IHandler(&genericHandler)
	return handler, nil
}
