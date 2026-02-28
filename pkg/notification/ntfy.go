package notification

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/TimothyYe/godns/internal/settings"
	log "github.com/sirupsen/logrus"
)

const defaultNtfyServer = "https://ntfy.sh"

type NtfyNotification struct {
	conf *settings.Settings
}

func NewNtfyNotification(conf *settings.Settings) INotification {
	return &NtfyNotification{conf: conf}
}

func (n *NtfyNotification) Send(domain, currentIP string) error {
	config := n.conf.Notify.Ntfy
	if config.Topic == "" {
		return errors.New("ntfy topic cannot be empty")
	}

	server := config.Server
	if server == "" {
		server = defaultNtfyServer
	}

	tpl := config.MsgTemplate
	if tpl == "" {
		tpl = "IP address of {{ .Domain }} updated to {{ .CurrentIP }}"
	}

	msg := buildTemplate(currentIP, domain, tpl)
	url := fmt.Sprintf("%s/%s", strings.TrimRight(server, "/"), config.Topic)

	req, err := http.NewRequest("POST", url, strings.NewReader(msg))
	if err != nil {
		return err
	}

	req.Header.Set("Title", "GoDNS Notification")

	if config.Priority != "" {
		req.Header.Set("Priority", config.Priority)
	}

	if config.Tags != "" {
		req.Header.Set("Tags", config.Tags)
	}

	if config.Icon != "" {
		req.Header.Set("Icon", config.Icon)
	}

	// Auth: Bearer token takes precedence over basic auth
	if config.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.Token))
	} else if config.User != "" && config.Password != "" {
		req.SetBasicAuth(config.User, config.Password)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send ntfy notification: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Debug("ntfy response Status: ", resp.Status)
	log.Debug("ntfy response Body: ", string(body))

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("ntfy notification failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
