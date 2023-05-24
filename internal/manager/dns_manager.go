package manager

import (
	"context"
	"os"

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

func (manager *DNSManager) SetConfiguration(conf *settings.Settings) *DNSManager {
	manager.configuration = conf
	return manager
}

func (manager *DNSManager) Build() error {
	log.Infof("Creating DNS handler with provider: %s", manager.configuration.Provider)
	dnsProvider, err := provider.GetProvider(manager.configuration)
	if err != nil {
		return err
	}

	manager.handler = &handler.Handler{}
	manager.handler.SetConfiguration(manager.configuration)
	manager.handler.SetProvider(dnsProvider)

	ctx, cancel := context.WithCancel(context.Background())
	manager.ctx = ctx
	manager.cancel = cancel
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
