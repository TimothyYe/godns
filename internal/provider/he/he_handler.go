package he

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/TimothyYe/godns/internal/settings"
	"github.com/TimothyYe/godns/internal/utils"
	log "github.com/sirupsen/logrus"
)

const (
	// URL the API address for he.net.
	URL = "https://dyn.dns.he.net/nic/update"
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

// updateIP update subdomain with current IP.
func (provider *DNSProvider) updateIP(domain, subDomain, currentIP string) error {
	values := url.Values{}

	if subDomain != utils.RootDomain {
		values.Add("hostname", fmt.Sprintf("%s.%s", subDomain, domain))
	} else {
		values.Add("hostname", domain)
	}
	values.Add("password", provider.configuration.Password)
	values.Add("myip", currentIP)

	client := utils.GetHTTPClient(provider.configuration)

	req, _ := http.NewRequest("POST", URL, strings.NewReader(values.Encode()))
	resp, err := client.Do(req)

	if err != nil {
		log.Error("Request error:", err)
		return err
	}

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusOK {
		log.Infof("Update IP success: %s", string(body))
	} else {
		log.Infof("Update IP failed: %s", string(body))
	}

	return nil
}
