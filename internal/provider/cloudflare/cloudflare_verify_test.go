package cloudflare

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/TimothyYe/godns/internal/settings"
)

type apiLog struct {
	mu      sync.Mutex
	updated []string // record IDs that received PUT
	deleted []string // record IDs that received DELETE
	created int
}

func newMockCloudflare(t *testing.T, records []DNSRecord, calls *apiLog) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()

	mux.HandleFunc("/zones", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(ZoneResponse{
			Zones:   []Zone{{ID: "zone1", Name: "example.com"}},
			Success: true,
		})
	})

	mux.HandleFunc("/zones/zone1/dns_records", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_ = json.NewEncoder(w).Encode(DNSRecordResponse{Records: records, Success: true})
		case http.MethodPost:
			calls.mu.Lock()
			calls.created++
			calls.mu.Unlock()
			_ = json.NewEncoder(w).Encode(DNSRecordUpdateResponse{Success: true})
		}
	})

	mux.HandleFunc("/zones/zone1/dns_records/", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/zones/zone1/dns_records/"):]
		calls.mu.Lock()
		switch r.Method {
		case http.MethodPut:
			calls.updated = append(calls.updated, id)
		case http.MethodDelete:
			calls.deleted = append(calls.deleted, id)
		}
		calls.mu.Unlock()
		_ = json.NewEncoder(w).Encode(DNSRecordUpdateResponse{Success: true})
	})

	return httptest.NewServer(mux)
}

func newTestProvider(apiURL string) *DNSProvider {
	return &DNSProvider{
		API: apiURL,
		configuration: &settings.Settings{
			IPType:     "IPv6",
			LoginToken: "test-token",
			Domains: []settings.Domain{
				{DomainName: "example.com", SubDomains: []string{"sub", "@"}},
			},
		},
	}
}

// Scenario from issue #268: flapping IPv6 left multiple AAAA records behind.
func TestDedupMultipleStaleRecords(t *testing.T) {
	calls := &apiLog{}
	records := []DNSRecord{
		{ID: "rec1", Name: "sub.example.com", IP: "fd00::old1", Type: "AAAA"},
		{ID: "rec2", Name: "sub.example.com", IP: "fd00::current", Type: "AAAA"},
		{ID: "rec3", Name: "sub.example.com", IP: "fd00::old2", Type: "AAAA"},
		{ID: "rec9", Name: "other.example.com", IP: "fd00::other", Type: "AAAA"}, // untracked
	}
	srv := newMockCloudflare(t, records, calls)
	defer srv.Close()

	provider := newTestProvider(srv.URL)
	if err := provider.UpdateIP("example.com", "sub", "fd00::current"); err != nil {
		t.Fatalf("UpdateIP failed: %v", err)
	}

	if len(calls.updated) != 1 || calls.updated[0] != "rec1" {
		t.Errorf("expected exactly rec1 updated, got %v", calls.updated)
	}
	if len(calls.deleted) != 2 {
		t.Errorf("expected 2 deletions, got %v", calls.deleted)
	}
	for _, id := range calls.deleted {
		if id == "rec9" {
			t.Errorf("untracked record rec9 must not be deleted")
		}
		if id == "rec1" {
			t.Errorf("kept record rec1 must not be deleted")
		}
	}
	if calls.created != 0 {
		t.Errorf("expected no creation, got %d", calls.created)
	}
}

// Single record with matching IP: nothing should change.
func TestSingleRecordUpToDate(t *testing.T) {
	calls := &apiLog{}
	records := []DNSRecord{
		{ID: "rec1", Name: "sub.example.com", IP: "fd00::current", Type: "AAAA"},
	}
	srv := newMockCloudflare(t, records, calls)
	defer srv.Close()

	provider := newTestProvider(srv.URL)
	if err := provider.UpdateIP("example.com", "sub", "fd00::current"); err != nil {
		t.Fatalf("UpdateIP failed: %v", err)
	}

	if len(calls.updated) != 0 || len(calls.deleted) != 0 || calls.created != 0 {
		t.Errorf("expected no API mutations, got updated=%v deleted=%v created=%d",
			calls.updated, calls.deleted, calls.created)
	}
}

// Root domain (@) records must be matched and deduplicated too.
func TestRootDomainDedup(t *testing.T) {
	calls := &apiLog{}
	records := []DNSRecord{
		{ID: "rec1", Name: "example.com", IP: "fd00::old1", Type: "AAAA"},
		{ID: "rec2", Name: "example.com", IP: "fd00::old2", Type: "AAAA"},
	}
	srv := newMockCloudflare(t, records, calls)
	defer srv.Close()

	provider := newTestProvider(srv.URL)
	if err := provider.UpdateIP("example.com", "@", "fd00::current"); err != nil {
		t.Fatalf("UpdateIP failed: %v", err)
	}

	if len(calls.updated) != 1 || calls.updated[0] != "rec1" {
		t.Errorf("expected rec1 updated, got %v", calls.updated)
	}
	if len(calls.deleted) != 1 || calls.deleted[0] != "rec2" {
		t.Errorf("expected rec2 deleted, got %v", calls.deleted)
	}
}

// No matching record: must fall through to creation.
func TestCreateWhenMissing(t *testing.T) {
	calls := &apiLog{}
	srv := newMockCloudflare(t, []DNSRecord{}, calls)
	defer srv.Close()

	provider := newTestProvider(srv.URL)
	if err := provider.UpdateIP("example.com", "sub", "fd00::current"); err != nil {
		t.Fatalf("UpdateIP failed: %v", err)
	}

	if calls.created != 1 {
		t.Errorf("expected 1 creation, got %d", calls.created)
	}
	if len(calls.updated) != 0 || len(calls.deleted) != 0 {
		t.Errorf("expected no update/delete, got updated=%v deleted=%v", calls.updated, calls.deleted)
	}
}

// Old code matched via strings.Contains: subdomain "sub" would also hit
// "mysub.example.com". Verify the new exact matching does not.
func TestNoSubstringFalseMatch(t *testing.T) {
	calls := &apiLog{}
	records := []DNSRecord{
		{ID: "recX", Name: "mysub.example.com", IP: "fd00::x", Type: "AAAA"},
	}
	srv := newMockCloudflare(t, records, calls)
	defer srv.Close()

	provider := newTestProvider(srv.URL)
	if err := provider.UpdateIP("example.com", "sub", "fd00::current"); err != nil {
		t.Fatalf("UpdateIP failed: %v", err)
	}

	if len(calls.updated) != 0 || len(calls.deleted) != 0 {
		t.Errorf("mysub.example.com must not be touched, got updated=%v deleted=%v",
			calls.updated, calls.deleted)
	}
	if calls.created != 1 {
		t.Errorf("expected sub.example.com to be created, got %d creations", calls.created)
	}
}
