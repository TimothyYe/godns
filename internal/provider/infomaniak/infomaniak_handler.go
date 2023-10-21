package infomaniak

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
	// URL the API address for Infomaniak.
	URL = "https://%s:%s@infomaniak.com/nic/update?hostname=%s.%s&myip=%s"
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
	client := utils.GetHTTPClient(provider.configuration)
	resp, err := client.Get(fmt.Sprintf(URL,
		provider.configuration.Email,
		provider.configuration.Password,
		subDomain,
		domain,
		currentIP))

	if err != nil {
		// handle error
		log.Error("Failed to update sub domain:", subDomain)
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error(err)
		}
	}(resp.Body)

	if err != nil {
		log.Error("Err:", err.Error())
		return err
	}

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Errorf("Update IP failed: %s", string(body))
		return fmt.Errorf("update IP failed: %s", string(body))
	}

	if strings.Contains(string(body), "good") {
		log.Infof("Update IP success: %s", string(body))
	} else if strings.Contains(string(body), "nochg") {
		log.Infof("IP not changed: %s", string(body))
	} else {
		return fmt.Errorf("update IP failed: %s", string(body))
	}

	return nil
}
