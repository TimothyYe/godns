package hetzner

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"path"
	"sync"
	"testing"

	"github.com/TimothyYe/godns/internal/settings"
	"github.com/TimothyYe/godns/internal/utils"
)

type apiLog struct {
	mu      sync.Mutex
	updated []Record // records that received a PUT.
}

func newMockHetzner(t *testing.T, calls *apiLog) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()

	mux.HandleFunc("/zones", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, `{"zones":[{"id":"zone1","name":"example.com"}]}`)
	})

	// sub has both an A and an AAAA record, the usual dual stack setup.
	mux.HandleFunc("/records", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, `{"records":[
			{"id":"recA","zone_id":"zone1","name":"sub","type":"A","value":"203.0.113.1","ttl":3600},
			{"id":"recAAAA","zone_id":"zone1","name":"sub","type":"AAAA","value":"2001:db8::1","ttl":3600}
		]}`)
	})

	mux.HandleFunc("/records/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var record Record
		if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
			t.Errorf("failed to decode PUT body for %s: %v", path.Base(r.URL.Path), err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		calls.mu.Lock()
		calls.updated = append(calls.updated, record)
		calls.mu.Unlock()
	})

	return httptest.NewServer(mux)
}

func newTestProvider(t *testing.T, apiURL, ipType string) *DNSProvider {
	t.Helper()
	provider := &DNSProvider{}
	provider.Init(&settings.Settings{
		IPType:     ipType,
		LoginToken: "test-token",
		Domains: []settings.Domain{
			{DomainName: "example.com", SubDomains: []string{"sub"}},
		},
	})
	provider.API = apiURL + "/"

	return provider
}

// ip_type is documented as "IPv4" or "IPv6" and every sample config writes it
// that way, while the web UI saves it as "IPV4" or "IPV6". Whichever spelling
// reaches the provider, an IPv6 configuration has to update the AAAA record.
func TestUpdateIPPicksRecordTypeForConfiguredIPType(t *testing.T) {
	tests := []struct {
		name       string
		ipType     string
		newIP      string
		wantRecord string
		wantType   string
	}{
		{name: "IPv6", ipType: "IPv6", newIP: "2001:db8::2", wantRecord: "recAAAA", wantType: utils.IPTypeAAAA},
		{name: "IPV6", ipType: "IPV6", newIP: "2001:db8::2", wantRecord: "recAAAA", wantType: utils.IPTypeAAAA},
		{name: "lowercase ipv6", ipType: "ipv6", newIP: "2001:db8::2", wantRecord: "recAAAA", wantType: utils.IPTypeAAAA},
		{name: "IPv4", ipType: "IPv4", newIP: "203.0.113.2", wantRecord: "recA", wantType: utils.IPTypeA},
		{name: "unset", ipType: "", newIP: "203.0.113.2", wantRecord: "recA", wantType: utils.IPTypeA},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			calls := &apiLog{}
			srv := newMockHetzner(t, calls)
			defer srv.Close()

			provider := newTestProvider(t, srv.URL, tc.ipType)
			if err := provider.UpdateIP("example.com", "sub", tc.newIP); err != nil {
				t.Fatalf("UpdateIP failed: %v", err)
			}

			if len(calls.updated) != 1 {
				t.Fatalf("ip_type %q: expected 1 record update, got %d", tc.ipType, len(calls.updated))
			}

			got := calls.updated[0]
			if got.ID != tc.wantRecord || got.Type != tc.wantType {
				t.Errorf("ip_type %q: updated record %s of type %s, want %s of type %s",
					tc.ipType, got.ID, got.Type, tc.wantRecord, tc.wantType)
			}
			if got.Value != tc.newIP {
				t.Errorf("ip_type %q: record value is %s, want %s", tc.ipType, got.Value, tc.newIP)
			}
		})
	}
}
