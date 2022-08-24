package lib

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"sync"
	"text/template"

	"github.com/TimothyYe/godns/internal/settings"
	"github.com/TimothyYe/godns/internal/utils"
	log "github.com/sirupsen/logrus"
)

type Webhook struct {
	conf   *settings.Settings
	client *http.Client
}

var (
	instance *Webhook
	once     sync.Once
)

func GetWebhook(conf *settings.Settings) *Webhook {
	once.Do(func() {
		instance = &Webhook{
			conf:   conf,
			client: utils.GetHTTPClient(conf),
		}
	})

	return instance
}

func (w *Webhook) Execute(domain, currentIP string) error {
	if w.conf.Webhook.URL == "" {
		log.Debug("Webhook URL is empty, skip sending notification")
		return nil
	}

	// set request method
	method := http.MethodGet
	if w.conf.Webhook.RequestBody != "" {
		method = http.MethodPost
	}

	reqURL, reqBody := "", ""
	var err error
	// send HTTP get request
	if method == http.MethodGet {
		reqURL, err = w.buildReqURL(domain, currentIP, w.conf.IPType)
		if err != nil {
			return err
		}
	} else {
		reqURL = w.conf.Webhook.URL
		reqBody, err = w.buildReqBody(domain, currentIP, w.conf.IPType)
		if err != nil {
			return err
		}
	}

	var req *http.Request
	req, err = http.NewRequest(method, reqURL, strings.NewReader(reqBody))
	if err != nil {
		log.Error("Failed to create request:", err)
		return err
	}

	if method == http.MethodPost {
		req.Header.Add("Content-Type", "application/json")
	}

	resp, err := w.client.Do(req)
	if err != nil {
		log.Error("Failed to send request:", err)
		return err
	}

	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Failed to read response body:", err)
		return err
	}

	log.Debugf("Webhook response: %s", string(content))
	return nil
}

func (w *Webhook) buildReqURL(domain, currentIP, ipType string) (string, error) {
	t := template.New("req template")
	if _, err := t.Parse(w.conf.Webhook.URL); err != nil {
		log.Error("Failed to parse template:", err)
		return "", err
	}

	data := struct {
		CurrentIP string
		Domain    string
		IPType    string
	}{
		currentIP,
		domain,
		ipType,
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		log.Error(err)
		return "", err
	}

	return tpl.String(), nil
}

func (w *Webhook) buildReqBody(domain, currentIP, ipType string) (string, error) {
	t := template.New("reqBody template")
	if _, err := t.Parse(w.conf.Webhook.RequestBody); err != nil {
		log.Error("Failed to parse template:", err)
		return "", err
	}

	data := struct {
		CurrentIP string
		Domain    string
		IPType    string
	}{
		currentIP,
		domain,
		ipType,
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		log.Error(err)
		return "", err
	}

	return tpl.String(), nil
}
