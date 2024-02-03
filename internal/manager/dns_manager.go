package manager

import (
	"context"
	"os"
	"sync"

	"github.com/TimothyYe/godns/internal/handler"
	"github.com/TimothyYe/godns/internal/provider"
	"github.com/TimothyYe/godns/internal/settings"
	log "github.com/sirupsen/logrus"
)

type DNSManager struct {
	configuration *settings.Settings
	handler       *handler.Handler
	provider      provider.IDNSProvider
	ctx           context.Context
	cancel        context.CancelFunc
}

var (
	managerInstance *DNSManager
	managerOnce     sync.Once
)

func GetDNSManager(conf *settings.Settings) *DNSManager {
	managerOnce.Do(func() {
		managerInstance = &DNSManager{}
		managerInstance.setConfiguration(conf)
		if err := managerInstance.initManager(); err != nil {
			log.Fatalf("Error during DNS manager initialization: %s", err)
		}
	})

	return managerInstance
}

func (manager *DNSManager) setConfiguration(conf *settings.Settings) *DNSManager {
	manager.configuration = conf
	return manager
}

func (manager *DNSManager) initManager() error {
	log.Infof("Creating DNS handler with provider: %s", manager.configuration.Provider)
	dnsProvider, err := provider.GetProvider(manager.configuration)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	manager.ctx = ctx
	manager.cancel = cancel

	manager.provider = dnsProvider
	manager.handler = &handler.Handler{}
	manager.handler.SetContext(manager.ctx)
	manager.handler.SetConfiguration(manager.configuration)
	manager.handler.SetProvider(manager.provider)

	return nil
}

func (manager *DNSManager) Run() {
	for _, domain := range manager.configuration.Domains {
		domain := domain

		if manager.configuration.RunOnce {
			err := manager.handler.UpdateIP(&domain)
			if err != nil {
				log.Error("Error during execution:", err)
				os.Exit(1)
			}
		} else {
			// pass the context to the goroutine
			go manager.handler.LoopUpdateIP(manager.ctx, &domain)
		}
	}

	if manager.configuration.RunOnce {
		os.Exit(0)
	}
}

func (manager *DNSManager) Stop() {
	manager.cancel()
}
