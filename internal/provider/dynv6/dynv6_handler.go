package dynv6

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/TimothyYe/godns/internal/settings"
	"github.com/TimothyYe/godns/internal/utils"
	log "github.com/sirupsen/logrus"
)

const (
	// URL the API address for Duck DNS.
	URL = "https://dynv6.com/api/update?hostname=%s&token=%s&%s"
)

// DNSProvider struct.
type DNSProvider struct {
	configuration *settings.Settings
}

// Init passes DNS settings and store it to provider instance.
func (provider *DNSProvider) Init(conf *settings.Settings) {
	provider.configuration = conf
}

func (provider *DNSProvider) UpdateIP(domainName, subdomainName, ip string) error {
	hostname := subdomainName + "." + domainName
	client := utils.GetHTTPClient(provider.configuration)
	return provider.update(client, hostname, ip)
}

func (provider *DNSProvider) update(client *http.Client, hostname string, currentIP string) error {
	var ip string
	if strings.ToUpper(provider.configuration.IPType) == utils.IPV4 {
		ip = fmt.Sprintf("ipv4=%s", currentIP)
	} else if strings.ToUpper(provider.configuration.IPType) == utils.IPV6 {
		ip = fmt.Sprintf("ipv6=%s", currentIP)
	}

	// update IP with HTTP GET request
	url := fmt.Sprintf(URL, hostname, provider.configuration.LoginToken, ip)
	log.Debug("Update url: ", url)
	resp, err := client.Get(url)
	if err != nil {
		// handle error
		log.Errorf("Cannot send request: %s", err.Error())
		return fmt.Errorf("cannot send request: %s", err.Error())
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error(err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to receive response: %w", err)
	} else if !strings.Contains(string(body), "addresses updated") {
		return fmt.Errorf("service rejected update: %s", string(body))
	}

	return nil
}
