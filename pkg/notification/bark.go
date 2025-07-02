package notification

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/TimothyYe/godns/internal/settings"
	log "github.com/sirupsen/logrus"
)

type BarkNotification struct {
	conf *settings.Settings
}

type BarkParams struct {
	Title      string   `json:"title"`
	Subtitle   string   `json:"subtitle"`
	Body       string   `json:"body"`
	DeviceKeys []string `json:"device_keys"`
	Level      string   `json:"level,omitempty"`
	Volume     int      `json:"volume,omitempty"`
	Badge      int      `json:"badge,omitempty"`
	Call       int      `json:"call,omitempty"`
	AutoCopy   int      `json:"autoCopy,omitempty"`
	Copy       string   `json:"copy,omitempty"`
	Sound      string   `json:"sound,omitempty"`
	Icon       string   `json:"icon,omitempty"`
	Group      string   `json:"group,omitempty"`
	IsArchive  int      `json:"is_archive,omitempty"`
	URL        string   `json:"url,omitempty"`
	Image      string   `json:"image,omitempty"`
	Action     string   `json:"action,omitempty"`
}

func NewBarkNotification(conf *settings.Settings) INotification {
	return &BarkNotification{conf: conf}
}

func (n *BarkNotification) Send(domain, currentIP string) error {
	if n.conf.Notify.Bark.DeviceKeys == "" {
		return errors.New("bark device keys cannot be empty")
	}

	if n.conf.Notify.Bark.Server == "" {
		n.conf.Notify.Bark.Server = "https://api.day.app"
	}

	params := &BarkParams{
		IsArchive: 1,
		Action:    "none",
	}
	if n.conf.Notify.Bark.Params != "" {
		if err := json.Unmarshal([]byte(n.conf.Notify.Bark.Params), params); err != nil {
			return err
		}
	}

	// set default params
	params.Title = "GoDNS Notification"
	params.Subtitle = "{{ .Domain }}"
	params.Body = "[{{ .CurrentIP }}]"
	params.DeviceKeys = strings.Split(n.conf.Notify.Bark.DeviceKeys, ",")
	// override title
	if title := n.conf.Notify.Bark.Title; title != "" {
		params.Title = title
	}
	// override subtitle
	if subtitle := n.conf.Notify.Bark.Subtitle; subtitle != "" {
		params.Subtitle = subtitle
	}
	// override body
	if body := n.conf.Notify.Bark.Body; body != "" {
		params.Body = body
	}

	tpl, err := json.Marshal(params)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/push", n.conf.Notify.Bark.Server)

	re := regexp.MustCompile(`\[(.*?)\]\s+of\s+(.*?)$`)
	matches := re.FindStringSubmatch(domain)
	if len(matches) != 3 {
		return errors.New("invalid format: expected '[...] of ...'")
	}

	rootDomain := strings.TrimSpace(matches[2])
	// extract subdomains
	subDomains := strings.Split(matches[1], ",")
	for i := range subDomains {
		domainName := fmt.Sprintf("%s.%s", strings.TrimSpace(subDomains[i]), rootDomain)
		data := buildTemplate(currentIP, domainName, string(tpl))
		body := bytes.NewBuffer([]byte(data))
		_ = n.sendJSON(url, body)
	}

	return nil
}

func (n *BarkNotification) sendJSON(url string, body io.Reader) error {
	// Create client
	client := &http.Client{}
	// Create request
	req, _ := http.NewRequest("POST", url, body)
	// Headers
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	// Fetch Request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failure : ", err)
		return err
	}

	// Read Response Body
	respBody, _ := io.ReadAll(resp.Body)
	// Display Results
	log.Debug("response Status : ", resp.Status)
	log.Debug("response Body : ", string(respBody))
	return nil
}
