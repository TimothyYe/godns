package noip

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
	// URL the API address for NoIP.
	URL = "https://%s:%s@dynupdate.no-ip.com/nic/update?hostname=%s&%s"
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
	return provider.update(client, hostname, subdomainName, ip)
}

func (provider *DNSProvider) update(client *http.Client, hostname, subdomain string, currentIP string) error {
	var ip string
	if strings.ToUpper(provider.configuration.IPType) == utils.IPV4 {
		ip = fmt.Sprintf("myip=%s", currentIP)
	} else if strings.ToUpper(provider.configuration.IPType) == utils.IPV6 {
		ip = fmt.Sprintf("myipv6=%s", currentIP)
	}

	req, _ := http.NewRequest("GET", fmt.Sprintf(
		URL,
		provider.configuration.Email,
		provider.configuration.Password,
		hostname,
		ip), nil)

	if provider.configuration.UserAgent != "" {
		req.Header.Add("User-Agent", provider.configuration.UserAgent)
	}

	// update IP with HTTP GET request
	resp, err := client.Do(req)
	if err != nil {
		// handle error
		log.Error("Failed to update sub domain:", subdomain)
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error(err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil || !strings.Contains(string(body), "good") {
		log.Error("Failed to update the IP", err)
		return err
	}

	log.Infof("IP updated to: %s", currentIP)
	return nil
}
