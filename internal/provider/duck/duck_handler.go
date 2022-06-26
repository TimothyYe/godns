package duck

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/TimothyYe/godns/internal/settings"
	"github.com/TimothyYe/godns/internal/utils"
	log "github.com/sirupsen/logrus"
)

const (
	// URL the API address for Duck DNS.
	URL = "https://www.duckdns.org/update?domains=%s&token=%s&%s"
)

// DNSProvider struct.
type DNSProvider struct {
	configuration *settings.Settings
}

// Init passes DNS settings and store it to the provider instance.
func (provider *DNSProvider) Init(conf *settings.Settings) {
	provider.configuration = conf
}

func (provider *DNSProvider) UpdateIP(domainName, subdomainName, ip string) error {
	return provider.updateIP(domainName, subdomainName, ip)
}

func (provider *DNSProvider) updateIP(domainName, subdomainName, currentIP string) error {
	var ip string

	if strings.ToUpper(provider.configuration.IPType) == utils.IPV4 {
		ip = fmt.Sprintf("ip=%s", currentIP)
	} else if strings.ToUpper(provider.configuration.IPType) == utils.IPV6 {
		ip = fmt.Sprintf("ipv6=%s", currentIP)
	}

	client := utils.GetHTTPClient(provider.configuration, provider.configuration.UseProxy)

	// update IP with HTTP GET request
	resp, err := client.Get(fmt.Sprintf(URL, subdomainName, provider.configuration.LoginToken, ip))
	if err != nil {
		// handle error
		log.Errorf("Failed to update sub domain: %s.%s", domainName, subdomainName)
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error(err)
		}
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil || string(body) != "OK" {
		log.Error("Failed to update the IP:", err)
		return err
	}

	log.Infof("IP updated to: %s", ip)
	return nil
}
