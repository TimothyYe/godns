package provider

import (
	"fmt"

	"github.com/TimothyYe/godns/internal/provider/alidns"
	"github.com/TimothyYe/godns/internal/provider/cloudflare"
	"github.com/TimothyYe/godns/internal/provider/dnspod"
	"github.com/TimothyYe/godns/internal/provider/dreamhost"
	"github.com/TimothyYe/godns/internal/provider/duck"
	"github.com/TimothyYe/godns/internal/provider/dynu"
	"github.com/TimothyYe/godns/internal/provider/dynv6"
	"github.com/TimothyYe/godns/internal/provider/google"
	"github.com/TimothyYe/godns/internal/provider/he"
	"github.com/TimothyYe/godns/internal/provider/hetzner"
	"github.com/TimothyYe/godns/internal/provider/infomaniak"
	"github.com/TimothyYe/godns/internal/provider/linode"
	"github.com/TimothyYe/godns/internal/provider/loopiase"
	"github.com/TimothyYe/godns/internal/provider/noip"
	"github.com/TimothyYe/godns/internal/provider/ovh"
	"github.com/TimothyYe/godns/internal/provider/scaleway"
	"github.com/TimothyYe/godns/internal/provider/strato"
	"github.com/TimothyYe/godns/internal/settings"
	"github.com/TimothyYe/godns/internal/utils"
)

func GetProvider(conf *settings.Settings) (IDNSProvider, error) {
	var provider IDNSProvider

	switch conf.Provider {
	case utils.CLOUDFLARE:
		provider = &cloudflare.DNSProvider{}
	case utils.DNSPOD:
		provider = &dnspod.DNSProvider{}
	case utils.DREAMHOST:
		provider = &dreamhost.DNSProvider{}
	case utils.HE:
		provider = &he.DNSProvider{}
	case utils.ALIDNS:
		provider = &alidns.DNSProvider{}
	case utils.GOOGLE:
		provider = &google.DNSProvider{}
	case utils.DUCK:
		provider = &duck.DNSProvider{}
	case utils.NOIP:
		provider = &noip.DNSProvider{}
	case utils.SCALEWAY:
		provider = &scaleway.DNSProvider{}
	case utils.DYNV6:
		provider = &dynv6.DNSProvider{}
	case utils.LINODE:
		provider = &linode.DNSProvider{}
	case utils.STRATO:
		provider = &strato.DNSProvider{}
	case utils.LOOPIASE:
		provider = &loopiase.DNSProvider{}
	case utils.INFOMANIAK:
		provider = &infomaniak.DNSProvider{}
	case utils.HETZNER:
		provider = &hetzner.DNSProvider{}
	case utils.OVH:
		provider = &ovh.DNSProvider{}
	case utils.DYNU:
		provider = &dynu.DNSProvider{}
	default:
		return nil, fmt.Errorf("Unknown provider '%s'", conf.Provider)
	}

	provider.Init(conf)

	return provider, nil
}
