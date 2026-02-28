package notification

import (
	"sync"

	"github.com/TimothyYe/godns/internal/settings"

	log "github.com/sirupsen/logrus"
)

const (
	Email    = "email"
	Slack    = "slack"
	Telegram = "telegram"
	Discord  = "discord"
	Pushover = "pushover"
	Bark     = "bark"
	Ntfy     = "ntfy"
)

var (
	instance *notificationManager
	once     sync.Once
)

type INotification interface {
	Send(domain, currentIP string) error
}

type INotificationManager interface {
	Send(string, string)
}

type notificationManager struct {
	notifications map[string]INotification
}

func GetNotificationManager(conf *settings.Settings) INotificationManager {
	once.Do(func() {
		instance = &notificationManager{
			notifications: initNotifications(conf),
		}
	})

	return instance
}

func initNotifications(conf *settings.Settings) map[string]INotification {
	notificationMap := map[string]INotification{}

	if conf.Notify.Mail.Enabled {
		notificationMap[Email] = NewEmailNotification(conf)
	}

	if conf.Notify.Slack.Enabled {
		notificationMap[Slack] = NewSlackNotification(conf)
	}

	if conf.Notify.Telegram.Enabled {
		notificationMap[Telegram] = NewTelegramNotification(conf)
	}

	if conf.Notify.Pushover.Enabled {
		notificationMap[Pushover] = NewPushoverNotification(conf)
	}

	if conf.Notify.Discord.Enabled {
		notificationMap[Discord] = NewDiscordNotification(conf)
	}

	if conf.Notify.Bark.Enabled {
		notificationMap[Bark] = NewBarkNotification(conf)
	}

	if conf.Notify.Ntfy.Enabled {
		notificationMap[Ntfy] = NewNtfyNotification(conf)
	}

	return notificationMap
}

func (n *notificationManager) Send(domain, currentIP string) {
	for _, sender := range n.notifications {
		if err := sender.Send(domain, currentIP); err != nil {
			log.Error("Send notification with error:", err)
		}
	}
}
