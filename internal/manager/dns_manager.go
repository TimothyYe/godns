package manager

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/TimothyYe/godns/internal/handler"
	"github.com/TimothyYe/godns/internal/provider"
	"github.com/TimothyYe/godns/internal/settings"
	"github.com/TimothyYe/godns/internal/utils"
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
	configPath    string
}

var (
	managerInstance *DNSManager
	managerOnce     sync.Once
)

func getFileName(configPath string) string {
	// get the file name from the path
	// e.g. /etc/godns/config.json -> config.json
	return filepath.Base(configPath)
}

func (manager *DNSManager) setConfig(conf *settings.Settings) {
	manager.configuration = conf
}

func GetDNSManager(cfgPath string, conf *settings.Settings) *DNSManager {
	managerOnce.Do(func() {
		managerInstance = &DNSManager{}
		managerInstance.configPath = cfgPath
		managerInstance.configuration = conf
		if err := managerInstance.initManager(); err != nil {
			log.Fatalf("Error during DNS manager initialization: %s", err)
		}
	})

	return managerInstance
}

func (manager *DNSManager) startMonitor(ctx context.Context) {
	// Start listening for events.
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Debug("Shutting down the old file watcher...")
				return
			case event, ok := <-manager.watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					log.Debug("modified file:", event.Name)
					log.Debug("Reloading configuration...")
					// reload the configuration
					// read the file and update the configuration
					configFile := getFileName(manager.configPath)
					if event.Name == configFile {
						// Load settings from configs file
						newConfig := &settings.Settings{}
						if err := settings.LoadSettings(manager.configPath, newConfig); err != nil {
							log.Errorf("Failed to reload configuration: %s", err)
							continue
						}

						// validate the new configuration
						if err := utils.CheckSettings(newConfig); err != nil {
							log.Errorf("Failed to validate the new configuration: %s", err)
							continue
						}

						manager.setConfig(newConfig)
						manager.Restart()
					}
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
	if err := manager.watcher.Add(manager.configPath); err != nil {
		log.Fatal(err)
	}
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
	manager.handler.Init()

	// if RunOnce is true, we don't need to create a file watcher
	if !manager.configuration.RunOnce {
		// create a new file watcher
		log.Debug("Creating the new file watcher...")
		managerInstance.watcher, err = fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}

		// monitor the configuration file changes
		managerInstance.startMonitor(ctx)
	}
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

func (manager *DNSManager) Restart() {
	log.Info("Restarting DNS manager...")
	manager.Stop()

	// wait for the goroutines to exit
	time.Sleep(200 * time.Millisecond)

	// re-init the manager
	if err := manager.initManager(); err != nil {
		log.Fatalf("Error during DNS manager restarting: %s", err)
	}

	manager.Run()
	log.Info("DNS manager restarted successfully")
}
