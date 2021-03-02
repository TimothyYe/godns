package notify

import (
	"sync"

	"github.com/TimothyYe/godns"
)

var (
	instance *notifyManager
	once     sync.Once
)

type INotify interface {
	Send(conf *godns.Settings, domain, currentIP string) error
}

type notifyManager struct {
	notifications map[string]*INotify
}

func GetNotifyManager(conf *godns.Settings) *notifyManager {
	once.Do(func() {
		instance = &notifyManager{
			notifications: initNotifications(conf),
		}
	})

	return instance
}

func initNotifications(conf *godns.Settings) map[string]*INotify {
	return map[string]*INotify{}
}
