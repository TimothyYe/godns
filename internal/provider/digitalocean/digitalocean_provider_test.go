package digitalocean

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/TimothyYe/godns/internal/settings"
)

func TestDNSResponseToJSON(t *testing.T) {
	s := strings.NewReader(`
    {
        "domain_records": [
            {
                "id": 12345678,
                "type": "A",
                "name": "potato",
                "data": "127.0.0.1",
                "priority": null,
                "port": null,
                "ttl": 3600,
                "weight": null,
                "flags": null,
                "tag": null
            }
        ],
        "links": {},
        "meta": {
            "total": 1
        }
    }`)

	var resp DomainRecordsResponse
	err := json.NewDecoder(s).Decode(&resp)
	if err != nil {
		t.Error(err.Error())
	}
	if resp.Records[0].ID != 12345678 {
		t.Errorf("ID Error: %#v != 12345678 ", resp.Records[0].ID)
	}
	if resp.Records[0].Name != "potato" {
		t.Errorf("Name Error: %#v != potato", resp.Records[0].Name)
	}
}
func TestDNSUpdateResponseToJSON(t *testing.T) {
	s := strings.NewReader(`
    {
        "id": 12345678,
        "type": "A",
        "name": "@",
        "data": "127.0.0.1",
        "priority": null,
        "port": null,
        "ttl": 3600,
        "weight": null,
        "flags": null,
        "tag": null
    }`)

	var resp DNSRecord
	err := json.NewDecoder(s).Decode(&resp)
	if err != nil {
		t.Error(err.Error())
	}
	if resp.ID != 12345678 {
		t.Errorf("ID Error: %#v != 12345678 ", resp.ID)
	}
	if resp.Name != "@" {
		t.Errorf("Name Error: %#v != @", resp.Name)
	}
}

func TestRecordTracked(t *testing.T) {
	s := strings.NewReader(`
    {
        "domain_records": [
            {
                "id": 12345678,
                "type": "A",
                "name": "@",
                "data": "127.0.0.1",
                "priority": null,
                "port": null,
                "ttl": 3600,
                "weight": null,
                "flags": null,
                "tag": null
            },
            {
                "id": 12345678,
                "type": "A",
                "name": "swordfish",
                "data": "127.0.0.1",
                "priority": null,
                "port": null,
                "ttl": 3600,
                "weight": null,
                "flags": null,
                "tag": null
            },
            {
                "id": 12345678,
                "type": "A",
                "name": "www",
                "data": "127.0.0.1",
                "priority": null,
                "port": null,
                "ttl": 3600,
                "weight": null,
                "flags": null,
                "tag": null
            }
        ],
        "links": {},
        "meta": {
            "total": 3
        }
    }`)

	var resp DomainRecordsResponse
	err := json.NewDecoder(s).Decode(&resp)
	if err != nil {
		t.Error(err.Error())
	}
	var matchedDomains int
	domain := &settings.Domain{
		DomainName: "example.com",
		SubDomains: []string{"www", "@"},
	}

	for _, rec := range resp.Records {
		if recordTracked(domain, &rec) {
			t.Logf("Record founded: %+v", rec.Name)
			matchedDomains++
		}
	}
	if matchedDomains != 2 {
		t.Errorf("Unexpected amount of domains matched: %#v != 2", matchedDomains)
	}
}
