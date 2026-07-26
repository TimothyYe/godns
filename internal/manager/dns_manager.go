package manager

import (
	"context"
	"errors"
	"fmt"
	"net/http"
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
	provider    provider.IDNSProvider            // Legacy single provider (for backward compatibility)
	providers   map[string]provider.IDNSProvider // Multi-provider support
	ctx         context.Context
	cancel      context.CancelFunc
	watcher     *fsnotify.Watcher
	server      *server.Server
	configPath  string
	defaultAddr string
	// restartMu serializes Restart() — a single config save can fire multiple
	// fsnotify events, and overlapping restarts leave the manager in a
	// partially-initialized state.
	restartMu sync.Mutex
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
			SetConfigPath(manager.configPath).
			Build()

		srv := manager.server
		go func() {
			// Don't log.Fatalf here — Stop() / Restart() shut the server down
			// intentionally, and crashing the whole process from a background
			// goroutine on every restart would be wrong. Filter out the
			// expected shutdown signal and report anything else.
			if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Errorf("Web server stopped with error: %v", err)
			}
		}()
	} else {
		log.Info("Web panel is disabled")
	}
}

func (manager *DNSManager) startMonitor() {
	// Capture the current ctx and watcher in locals so this goroutine tracks
	// the lifecycle it was started with. Restart() reassigns these fields on
	// the manager — without local capture, an old monitor goroutine would
	// silently start reading from the new watcher's channels and double up
	// with the freshly-spawned monitor.
	ctx := manager.ctx
	watcher := manager.watcher

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Debug("Shutting down the old file watcher and the internal HTTP server...")
				return
			case event, ok := <-watcher.Events:
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
						// Restart() cancels ctx and closes watcher — exit so
						// the freshly-spawned monitor owns the new lifecycle.
						return
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Add a path. If this fails the live-reload feature is degraded for this
	// session, but it's not fatal — the rest of the app (DNS updates, web
	// panel) can run fine. The monitor goroutine will idle harmlessly with
	// no paths registered until shutdown closes ctx.
	if err := watcher.Add(manager.configPath); err != nil {
		log.Errorf("Failed to watch config file %s: %v (live config reload disabled)", manager.configPath, err)
	}
}

func (manager *DNSManager) initManager() error {
	ctx, cancel := context.WithCancel(context.Background())
	manager.ctx = ctx
	manager.cancel = cancel

	// Initialize providers based on configuration
	if manager.config.IsMultiProvider() {
		log.Info("Creating DNS handlers with multiple providers")
		providers, err := provider.GetProviders(manager.config)
		if err != nil {
			return fmt.Errorf("failed to initialize providers: %w", err)
		}
		manager.providers = providers

		// Log configured providers
		for providerName := range providers {
			log.Infof("Initialized provider: %s", providerName)
		}
	} else {
		log.Infof("Creating DNS handler with provider: %s", manager.config.Provider)
		dnsProvider, err := provider.GetProvider(manager.config)
		if err != nil {
			return err
		}
		manager.provider = dnsProvider
	}

	manager.handler = &handler.Handler{}
	manager.handler.SetContext(manager.ctx)
	manager.handler.SetConfiguration(manager.config)

	// Set provider(s) on handler
	if manager.config.IsMultiProvider() {
		manager.handler.SetProviders(manager.providers)
	} else {
		manager.handler.SetProvider(manager.provider)
	}

	manager.handler.Init()

	// if RunOnce is true, we don't need to create a file watcher and start the internal HTTP server
	if !manager.config.RunOnce {
		// create a new file watcher
		log.Debug("Creating the new file watcher...")
		var err error
		manager.watcher, err = fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}

		// monitor the configuration file changes
		manager.startMonitor()
		// start the internal HTTP server
		manager.startServer()
	}
	return nil
}

func (manager *DNSManager) Run() {
	if len(manager.config.Domains) == 0 {
		log.Info("No domain is configured, please check your configuration file")
		return
	}

	for _, domain := range manager.config.Domains {
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
	// Serialize restarts — fsnotify can fire multiple Write events for a
	// single editor save, and overlapping Stop()/initManager() calls leave
	// the manager with mismatched ctx/watcher/handler fields.
	manager.restartMu.Lock()
	defer manager.restartMu.Unlock()

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
