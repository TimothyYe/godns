package utils

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/TimothyYe/godns/internal/settings"

	log "github.com/sirupsen/logrus"

	"golang.org/x/net/proxy"
)

// GetHTTPClient creates the HTTP client and return it.
func GetHTTPClient(conf *settings.Settings) *http.Client {
	client := &http.Client{
		Timeout: time.Second * DefaultTimeout,
	}

	// Tune for the godns workload: a handful of long-lived destinations
	// (DNS provider APIs, IP-lookup services) hit on a recurring interval.
	// We want idle keep-alives across the loop iterations rather than a
	// fresh TLS handshake every refresh.
	httpTransport := &http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: conf.SkipSSLVerify},
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	if conf.UseProxy && conf.Socks5Proxy != "" {
		log.Debug("use socks5 proxy:" + conf.Socks5Proxy)
		dialer, err := proxy.SOCKS5("tcp", conf.Socks5Proxy, nil, proxy.Direct)
		if err != nil {
			log.Error("can't connect to the proxy:", err)
			return nil
		}

		httpTransport.DialContext = func(_ context.Context, network, address string) (net.Conn, error) {
			return dialer.Dial(network, address)
		}
	} else {
		httpTransport.DialContext = (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext
	}

	client.Transport = httpTransport
	return client
}
