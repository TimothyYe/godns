package notify

import (
	"log"
	"sync"

	"github.com/TimothyYe/godns"
)

const (
	Email    = "email"
	Slack    = "slack"
	Telegram = "telegram"
)

var (
	instance *notifyManager
	once     sync.Once
)

type INotify interface {
	Send(domain, currentIP string) error
}

type notifyManager struct {
	notifications map[string]INotify
}

func GetNotifyManager(conf *godns.Settings) *notifyManager {
	once.Do(func() {
		instance = &notifyManager{
			notifications: initNotifications(conf),
		}
	})

	return instance
}

func initNotifications(conf *godns.Settings) map[string]INotify {
	notifyMap := map[string]INotify{}

	if conf.Notify.Mail.Enabled {
		notifyMap[Email] = NewEmailNotify(conf)
	}

	if conf.Notify.Slack.Enabled {
		notifyMap[Slack] = NewSlackNotify(conf)
	}

	if conf.Notify.Telegram.Enabled {
		notifyMap[Telegram] = NewTelegramNotify(conf)
	}

	return notifyMap
}

func (n *notifyManager) Send(domain, currentIP string) {
	for _, sender := range n.notifications {
		err := sender.Send(domain, currentIP)
		log.Println("Send notification with error:", err.Error())
	}
}
