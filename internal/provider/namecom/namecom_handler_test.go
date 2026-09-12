package namecom

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/TimothyYe/godns/internal/settings"
	"github.com/TimothyYe/godns/internal/utils"
)

type apiCall struct {
	method string
	path   string
	body   recordBody
}

type mockAPI struct {
	mu    sync.Mutex
	calls []apiCall
	// pages holds the records returned for each page number, starting at 1.
	pages [][]Record
}

func (m *mockAPI) writes() []apiCall {
	m.mu.Lock()
	defer m.mu.Unlock()
	var out []apiCall
	for _, c := range m.calls {
		if c.method != http.MethodGet {
			out = append(out, c)
		}
	}
	return out
}

func newMockServer(t *testing.T, api *mockAPI) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, token, ok := r.BasicAuth()
		if !ok || user != "test-user" || token != "test-token" {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = io.WriteString(w, `{"message":"Unauthorized"}`)
			return
		}

		if !strings.HasPrefix(r.URL.Path, "/domains/example.com/records") {
			w.WriteHeader(http.StatusNotFound)
			_, _ = io.WriteString(w, `{"message":"Not Found"}`)
			return
		}

		call := apiCall{method: r.Method, path: r.URL.Path}
		if r.Method != http.MethodGet {
			if err := json.NewDecoder(r.Body).Decode(&call.body); err != nil {
				t.Errorf("failed to decode %s body: %v", r.Method, err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		api.mu.Lock()
		api.calls = append(api.calls, call)
		api.mu.Unlock()

		switch r.Method {
		case http.MethodGet:
			page, _ := strconv.Atoi(r.URL.Query().Get("page"))
			if page < 1 {
				page = 1
			}
			resp := ListRecordsResponse{}
			if page <= len(api.pages) {
				resp.Records = api.pages[page-1]
			}
			if page < len(api.pages) {
				resp.NextPage = page + 1
			}
			_ = json.NewEncoder(w).Encode(resp)
		case http.MethodPost:
			_, _ = fmt.Fprintf(w, `{"id":99,"domainName":"example.com","host":%q,"type":%q,"answer":%q,"ttl":%d}`,
				call.body.Host, call.body.Type, call.body.Answer, call.body.TTL)
		case http.MethodPut:
			_, _ = fmt.Fprintf(w, `{"id":%s,"domainName":"example.com","host":%q,"type":%q,"answer":%q,"ttl":%d}`,
				r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:], call.body.Host, call.body.Type, call.body.Answer, call.body.TTL)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
}

func newTestProvider(apiURL, ipType string) *DNSProvider {
	provider := &DNSProvider{}
	provider.Init(&settings.Settings{
		IPType:     ipType,
		Email:      "test-user",
		LoginToken: "test-token",
	})
	provider.API = apiURL
	return provider
}

func TestUpdateIP(t *testing.T) {
	records := []Record{
		{ID: 1, Host: "", Type: "A", Answer: "203.0.113.1", TTL: 300},
		{ID: 2, Host: "www", Type: "A", Answer: "203.0.113.1", TTL: 3600},
		{ID: 3, Host: "www", Type: "AAAA", Answer: "2001:db8::1", TTL: 300},
		{ID: 4, Host: "mail", Type: "MX", Answer: "mail.example.com", TTL: 300, Priority: 10},
	}

	tests := []struct {
		name       string
		ipType     string
		subdomain  string
		ip         string
		wantMethod string
		wantPath   string
		wantHost   string
		wantType   string
		wantTTL    int64
	}{
		{
			name: "update existing A record keeps its TTL", ipType: "IPv4", subdomain: "www", ip: "203.0.113.2",
			wantMethod: http.MethodPut, wantPath: "/domains/example.com/records/2", wantHost: "www", wantType: utils.IPTypeA, wantTTL: 3600,
		},
		{
			name: "update existing AAAA record", ipType: "IPv6", subdomain: "www", ip: "2001:db8::2",
			wantMethod: http.MethodPut, wantPath: "/domains/example.com/records/3", wantHost: "www", wantType: utils.IPTypeAAAA, wantTTL: 300,
		},
		{
			name: "web UI spelling of ip_type", ipType: "IPV6", subdomain: "www", ip: "2001:db8::2",
			wantMethod: http.MethodPut, wantPath: "/domains/example.com/records/3", wantHost: "www", wantType: utils.IPTypeAAAA, wantTTL: 300,
		},
		{
			name: "root domain matches empty host", ipType: "IPv4", subdomain: "@", ip: "203.0.113.2",
			wantMethod: http.MethodPut, wantPath: "/domains/example.com/records/1", wantHost: "", wantType: utils.IPTypeA, wantTTL: 300,
		},
		{
			name: "missing record is created", ipType: "IPv4", subdomain: "new", ip: "203.0.113.2",
			wantMethod: http.MethodPost, wantPath: "/domains/example.com/records", wantHost: "new", wantType: utils.IPTypeA, wantTTL: DefaultTTL,
		},
		{
			name: "missing AAAA for existing A host is created", ipType: "IPv6", subdomain: "mail", ip: "2001:db8::2",
			wantMethod: http.MethodPost, wantPath: "/domains/example.com/records", wantHost: "mail", wantType: utils.IPTypeAAAA, wantTTL: DefaultTTL,
		},
		{
			name: "unchanged IP makes no write", ipType: "IPv4", subdomain: "www", ip: "203.0.113.1",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			api := &mockAPI{pages: [][]Record{records}}
			srv := newMockServer(t, api)
			defer srv.Close()

			provider := newTestProvider(srv.URL, tc.ipType)
			if err := provider.UpdateIP("example.com", tc.subdomain, tc.ip); err != nil {
				t.Fatalf("UpdateIP failed: %v", err)
			}

			writes := api.writes()
			if tc.wantMethod == "" {
				if len(writes) != 0 {
					t.Fatalf("expected no write, got %+v", writes)
				}
				return
			}
			if len(writes) != 1 {
				t.Fatalf("expected 1 write, got %d: %+v", len(writes), writes)
			}
			got := writes[0]
			if got.method != tc.wantMethod || got.path != tc.wantPath {
				t.Errorf("got %s %s, want %s %s", got.method, got.path, tc.wantMethod, tc.wantPath)
			}
			if got.body.Host != tc.wantHost || got.body.Type != tc.wantType || got.body.Answer != tc.ip || got.body.TTL != tc.wantTTL {
				t.Errorf("got body %+v, want host=%q type=%q answer=%q ttl=%d",
					got.body, tc.wantHost, tc.wantType, tc.ip, tc.wantTTL)
			}
		})
	}
}

func TestListRecordsFollowsPagination(t *testing.T) {
	api := &mockAPI{pages: [][]Record{
		{{ID: 1, Host: "a", Type: "A", Answer: "203.0.113.1", TTL: 300}},
		{{ID: 2, Host: "b", Type: "A", Answer: "203.0.113.1", TTL: 300}},
		{{ID: 3, Host: "target", Type: "A", Answer: "203.0.113.1", TTL: 300}},
	}}
	srv := newMockServer(t, api)
	defer srv.Close()

	provider := newTestProvider(srv.URL, "IPv4")
	if err := provider.UpdateIP("example.com", "target", "203.0.113.2"); err != nil {
		t.Fatalf("UpdateIP failed: %v", err)
	}

	writes := api.writes()
	if len(writes) != 1 || writes[0].method != http.MethodPut || writes[0].path != "/domains/example.com/records/3" {
		t.Fatalf("expected PUT to record 3 found on the last page, got %+v", writes)
	}
}

func TestUpdateIPReportsAPIErrors(t *testing.T) {
	api := &mockAPI{}
	srv := newMockServer(t, api)
	defer srv.Close()

	provider := newTestProvider(srv.URL, "IPv4")
	provider.configuration.LoginToken = "wrong-token"

	err := provider.UpdateIP("example.com", "www", "203.0.113.2")
	if err == nil {
		t.Fatal("expected an error for rejected credentials")
	}
	if !strings.Contains(err.Error(), "401") || !strings.Contains(err.Error(), "Unauthorized") {
		t.Errorf("error should carry status and API message, got: %v", err)
	}
}
