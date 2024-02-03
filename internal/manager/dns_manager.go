package manager

import (
	"context"
	"os"
	"sync"

	"github.com/TimothyYe/godns/internal/handler"
	"github.com/TimothyYe/godns/internal/provider"
	"github.com/TimothyYe/godns/internal/settings"
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
)

type DNSManager struct {
	configuration *settings.Settings
	handler       *handler.Handler
	provider      provider.IDNSProvider
	ctx           context.Context
	cancel        context.CancelFunc
	watcher       *fsnotify.Watcher
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

		// create a new file watcher
		var err error
		managerInstance.watcher, err = fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}

		// monitor the configuration file changes
		managerInstance.startMonitor()
	})

	return managerInstance
}

func (manager *DNSManager) startMonitor() {
	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-manager.watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					log.Debug("modified file:", event.Name)
				}
			case err, ok := <-manager.watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Add a path.
	if err := manager.watcher.Add("./"); err != nil {
		log.Fatal(err)
	}
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
	manager.watcher.Close()
}
