package ionos

import (
	"fmt"
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
	mu             sync.Mutex
	requestedTypes []string // recordType query parameter of every record lookup.
	updated        []string // record IDs that received a PUT.
}

// sub.example.com has both an A and an AAAA record, the usual dual stack setup.
var zoneRecords = map[string]string{
	utils.IPTypeA:    `{"id":"recA","name":"sub.example.com","type":"A","content":"203.0.113.1","ttl":3600}`,
	utils.IPTypeAAAA: `{"id":"recAAAA","name":"sub.example.com","type":"AAAA","content":"2001:db8::1","ttl":3600}`,
}

func newMockIonos(t *testing.T, calls *apiLog) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()

	mux.HandleFunc("/zones", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, `[{"id":"zone1","name":"example.com","type":"NATIVE"}]`)
	})

	// The zone endpoint filters records by the recordName and recordType query
	// parameters, so a lookup for the wrong type comes back empty.
	mux.HandleFunc("/zones/zone1", func(w http.ResponseWriter, r *http.Request) {
		recordType := r.URL.Query().Get("recordType")

		calls.mu.Lock()
		calls.requestedTypes = append(calls.requestedTypes, recordType)
		calls.mu.Unlock()

		_, _ = fmt.Fprintf(w, `{"id":"zone1","name":"example.com","type":"NATIVE","records":[%s]}`,
			zoneRecords[recordType])
	})

	mux.HandleFunc("/zones/zone1/records/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		calls.mu.Lock()
		calls.updated = append(calls.updated, path.Base(r.URL.Path))
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
		wantType   string
		wantRecord string
	}{
		{name: "IPv6", ipType: "IPv6", newIP: "2001:db8::2", wantType: utils.IPTypeAAAA, wantRecord: "recAAAA"},
		{name: "IPV6", ipType: "IPV6", newIP: "2001:db8::2", wantType: utils.IPTypeAAAA, wantRecord: "recAAAA"},
		{name: "lowercase ipv6", ipType: "ipv6", newIP: "2001:db8::2", wantType: utils.IPTypeAAAA, wantRecord: "recAAAA"},
		{name: "AAAA", ipType: utils.IPTypeAAAA, newIP: "2001:db8::2", wantType: utils.IPTypeAAAA, wantRecord: "recAAAA"},
		{name: "IPv4", ipType: "IPv4", newIP: "203.0.113.2", wantType: utils.IPTypeA, wantRecord: "recA"},
		{name: "unset", ipType: "", newIP: "203.0.113.2", wantType: utils.IPTypeA, wantRecord: "recA"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			calls := &apiLog{}
			srv := newMockIonos(t, calls)
			defer srv.Close()

			provider := newTestProvider(t, srv.URL, tc.ipType)
			if err := provider.UpdateIP("example.com", "sub", tc.newIP); err != nil {
				t.Fatalf("UpdateIP failed: %v", err)
			}

			if len(calls.requestedTypes) != 1 || calls.requestedTypes[0] != tc.wantType {
				t.Errorf("ip_type %q: looked up record types %v, want [%s]",
					tc.ipType, calls.requestedTypes, tc.wantType)
			}
			if len(calls.updated) != 1 || calls.updated[0] != tc.wantRecord {
				t.Errorf("ip_type %q: updated records %v, want [%s]",
					tc.ipType, calls.updated, tc.wantRecord)
			}
		})
	}
}

// A record that already holds the current IP must not be written again.
func TestUpdateIPSkipsUnchangedRecord(t *testing.T) {
	calls := &apiLog{}
	srv := newMockIonos(t, calls)
	defer srv.Close()

	provider := newTestProvider(t, srv.URL, "IPv6")
	if err := provider.UpdateIP("example.com", "sub", "2001:db8::1"); err != nil {
		t.Fatalf("UpdateIP failed: %v", err)
	}

	if len(calls.updated) != 0 {
		t.Errorf("expected no record update, got %v", calls.updated)
	}
}
