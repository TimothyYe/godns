package lib

import (
	"net/http"
	"sync"

	"github.com/TimothyYe/godns/internal/utils"

	"github.com/TimothyYe/godns/internal/settings"
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
	return nil
}
