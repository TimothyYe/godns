package notify

import (
	"bytes"
	"log"
	"net/http"
	"text/template"

	"github.com/TimothyYe/godns"
	"golang.org/x/net/proxy"
)

// GetHttpClient creates the HTTP client and return it
func GetHttpClient(conf *godns.Settings, useProxy bool) *http.Client {
	client := &http.Client{}

	if useProxy && conf.Socks5Proxy != "" {
		log.Println("use socks5 proxy:" + conf.Socks5Proxy)
		dialer, err := proxy.SOCKS5("tcp", conf.Socks5Proxy, nil, proxy.Direct)
		if err != nil {
			log.Println("can't connect to the proxy:", err)
			return nil
		}

		httpTransport := &http.Transport{}
		client.Transport = httpTransport
		httpTransport.Dial = dialer.Dial
	}

	return client
}

func buildTemplate(currentIP, domain string, tplsrc string) string {
	t := template.New("notification template")
	if _, err := t.Parse(tplsrc); err != nil {
		log.Println("Failed to parse template")
		return ""
	}

	data := struct {
		CurrentIP string
		Domain    string
	}{
		currentIP,
		domain,
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		log.Println(err.Error())
		return ""
	}

	return tpl.String()
}
