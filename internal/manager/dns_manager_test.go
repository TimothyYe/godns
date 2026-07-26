package manager

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/TimothyYe/godns/internal/settings"
)

// writeTempConfig drops a minimal config file in a temp dir so the
// manager has a real path to hand to the fsnotify watcher. The
// contents don't matter — the file is never reloaded by these tests.
func writeTempConfig(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	if err := os.WriteFile(path, []byte(`{}`), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

// newTestManager builds a fully-initialized DNSManager around a fake
// configuration: a valid provider name (so factory.GetProvider succeeds),
// no domains (so Run is a no-op), and the web panel disabled (so no
// real listener is started). The caller is responsible for Stop().
func newTestManager(t *testing.T) *DNSManager {
	t.Helper()

	configPath := writeTempConfig(t)
	config := &settings.Settings{
		Provider:   "Cloudflare",
		Email:      "test@example.com",
		LoginToken: "test-token",
		Interval:   60,
		// No Domains and WebPanel disabled — Run becomes a no-op and
		// startServer skips listener startup.
	}

	m := &DNSManager{
		configPath:  configPath,
		config:      config,
		defaultAddr: ":9000", // matches startServer's "skip" condition
	}
	if err := m.initManager(); err != nil {
		t.Fatalf("initManager: %v", err)
	}
	return m
}

// TestRestart_BasicMechanics drives a single Restart cycle and asserts
// the manager's lifecycle fields are replaced consistently. The mutex
// is also exercised (single-threaded acquire/release).
func TestRestart_BasicMechanics(t *testing.T) {
	m := newTestManager(t)
	defer m.Stop()

	oldCtx := m.ctx
	oldWatcher := m.watcher

	m.Restart()

	if m.ctx == oldCtx {
		t.Error("expected ctx to be replaced after Restart")
	}
	if m.ctx.Err() != nil {
		t.Errorf("new ctx should not already be cancelled: %v", m.ctx.Err())
	}
	if m.watcher == oldWatcher {
		t.Error("expected watcher to be replaced after Restart")
	}
	if m.handler == nil {
		t.Error("handler is nil after Restart")
	}
}

// TestRestart_SerializedConcurrent fires several Restart calls in
// parallel and verifies (a) they all complete, and (b) the final
// manager state is internally consistent. The restartMu mutex is what
// makes this possible — without it, overlapping Stop/initManager runs
// would leave ctx/watcher/handler from different generations. Run this
// with -race to catch any field races.
func TestRestart_SerializedConcurrent(t *testing.T) {
	m := newTestManager(t)
	defer m.Stop()

	const goroutines = 4
	var wg sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			m.Restart()
		}()
	}

	finished := make(chan struct{})
	go func() {
		wg.Wait()
		close(finished)
	}()

	// Each Restart sleeps ~200ms; with serialization, total wall time
	// scales with goroutines. 10s is plenty even on a busy CI box.
	select {
	case <-finished:
	case <-time.After(10 * time.Second):
		t.Fatal("concurrent Restart calls did not finish within 10s")
	}

	if m.ctx == nil {
		t.Fatal("ctx is nil after concurrent restarts")
	}
	if m.ctx.Err() != nil {
		t.Errorf("ctx should be live after the final Restart: %v", m.ctx.Err())
	}
	if m.watcher == nil {
		t.Error("watcher is nil after concurrent restarts")
	}
	if m.handler == nil {
		t.Error("handler is nil after concurrent restarts")
	}
}
