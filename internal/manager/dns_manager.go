package manager

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/TimothyYe/godns/internal/handler"
	"github.com/TimothyYe/godns/internal/provider"
	"github.com/TimothyYe/godns/internal/server"
	"github.com/TimothyYe/godns/internal/settings"
	"github.com/TimothyYe/godns/internal/utils"
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
)

type DNSManager struct {
	config      *settings.Settings
	handler     *handler.Handler
	provider    provider.IDNSProvider
	ctx         context.Context
	cancel      context.CancelFunc
	watcher     *fsnotify.Watcher
	server      *server.Server
	configPath  string
	defaultAddr string
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

func GetDNSManager(cfgPath string, conf *settings.Settings, defaultAddr string) *DNSManager {
	managerOnce.Do(func() {
		managerInstance = &DNSManager{}
		managerInstance.configPath = cfgPath
		managerInstance.config = conf
		managerInstance.defaultAddr = defaultAddr
		if err := managerInstance.initManager(); err != nil {
			log.Fatalf("Error during DNS manager initialization: %s", err)
		}
	})

	return managerInstance
}

func (manager *DNSManager) startServer() {
	// start the internal HTTP server
	if (manager.config.WebPanel.Addr != "" || manager.defaultAddr != ":9000") && manager.config.WebPanel.Enabled {
		manager.server = &server.Server{}
		var addr string
		if manager.config.WebPanel.Addr != "" {
			addr = manager.config.WebPanel.Addr
		} else {
			addr = manager.defaultAddr
		}
		manager.server.
			SetAddress(addr).
			SetAuthInfo(manager.config.WebPanel.Username, manager.config.WebPanel.Password).
			SetConfig(manager.config).
			Build()

		go func() {
			if err := manager.server.Start(); err != nil {
				log.Fatalf("Failed to start the web server, error:%v", err)
			}
		}()
	} else {
		log.Info("Web panel is disabled")
	}
}

func (manager *DNSManager) startMonitor() {
	// Start listening for events.
	go func() {
		for {
			select {
			case <-manager.ctx.Done():
				log.Debug("Shutting down the old file watcher and the internal HTTP server...")
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

						manager.config = newConfig
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
	log.Infof("Creating DNS handler with provider: %s", manager.config.Provider)
	dnsProvider, err := provider.GetProvider(manager.config)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	manager.ctx = ctx
	manager.cancel = cancel

	manager.provider = dnsProvider
	manager.handler = &handler.Handler{}
	manager.handler.SetContext(manager.ctx)
	manager.handler.SetConfiguration(manager.config)
	manager.handler.SetProvider(manager.provider)
	manager.handler.Init()

	// if RunOnce is true, we don't need to create a file watcher and start the internal HTTP server
	if !manager.config.RunOnce {
		// create a new file watcher
		log.Debug("Creating the new file watcher...")
		managerInstance.watcher, err = fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}

		// monitor the configuration file changes
		managerInstance.startMonitor()
		// start the internal HTTP server
		managerInstance.startServer()
	}
	return nil
}

func (manager *DNSManager) Run() {
	for _, domain := range manager.config.Domains {
		domain := domain

		if manager.config.RunOnce {
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

	if manager.config.RunOnce {
		os.Exit(0)
	}
}

func (manager *DNSManager) Stop() {
	manager.cancel()
	// close the file watcher
	if manager.watcher != nil {
		manager.watcher.Close()
	}

	// stop the internal HTTP server
	if manager.server != nil {
		manager.server.Stop()
	}
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
