package lib

import (
	"bytes"
	"net/http"
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

	return nil
}

func (w *Webhook) buildReqURL(domain, currentIP, ipType string) (string, error) {
	t := template.New("notification template")
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
	t := template.New("notification template")
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
