package linode

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/proxy"
	"golang.org/x/oauth2"

	"github.com/TimothyYe/godns/internal/settings"
)

func CreateHTTPClient(conf *settings.Settings) (*http.Client, error) {
	transport := &http.Transport{}
	transport, err := applyProxy(conf.Socks5Proxy, transport)
	if err != nil {
		log.Infof("Error connecting to proxy: '%s'", err)
		log.Info("Continuing without proxy")
	}

	if conf.LoginToken == "" {
		return nil, errors.New("LoginToken cannot be an empty string")
	}
	roundTripper := addBearerAuth(conf.LoginToken, transport)

	httpClient := http.Client{
		Timeout:   time.Second * 10,
		Transport: roundTripper,
	}
	return &httpClient, nil
}

func applyProxy(proxyAddress string, transport *http.Transport) (*http.Transport, error) {
	if proxyAddress == "" {
		log.Debug("Skipping proxy: proxy address is empty string")
		return transport, nil
	}
	dialer, err := proxy.SOCKS5("tcp", proxyAddress, nil, proxy.Direct)
	if err != nil {
		return transport, err
	}
	log.Infof("Connected to proxy : %s", proxyAddress)

	dialContext := func(ctx context.Context, network, address string) (net.Conn, error) {
		return dialer.Dial(network, address)
	}

	transport.DialContext = dialContext
	return transport, nil
}

func addBearerAuth(accessToken string, transport http.RoundTripper) http.RoundTripper {
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	transportWithAuth := &oauth2.Transport{
		Source: tokenSource,
		Base:   transport,
	}
	log.Debug("Using OAuth / API token to connect to DNS service")
	return transportWithAuth
}
