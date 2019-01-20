package handler

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/TimothyYe/godns"
)

func TestResponseToJSON(t *testing.T) {
	s := strings.NewReader(`
    {
        "errors": [],
        "messages": [],
        "result": [
            {
                "id": "mk2b6fa491c12445a4376666a32429e1",
                "name": "example.com",
                "status": "active"
            }
        ],
        "result_info": {
            "count": 1,
            "page": 1,
            "per_page": 20,
            "total_count": 1,
            "total_pages": 1
        },
        "success": true
    }`)

	var resp ZoneResponse
	err := json.NewDecoder(s).Decode(&resp)
	if err != nil {
		t.Error(err.Error())
	}
	if resp.Success != true {
		t.Errorf("Success Error: %#v != true ", resp.Success)
	}
	if resp.Zones[0].ID != "mk2b6fa491c12445a4376666a32429e1" {
		t.Errorf("ID Error: %#v != mk2b6fa491c12445a4376666a32429e1 ", resp.Zones[0].ID)
	}
	if resp.Zones[0].Name != "example.com" {
		t.Errorf("Name Error: %#v != example.com", resp.Zones[0].Name)
	}
}

func TestDNSResponseToJSON(t *testing.T) {
	s := strings.NewReader(`
    {
        "errors": [],
        "messages": [],
        "result": [
            {
                "content": "127.0.0.1",
                "id": "F11cc63e02a42d38174b8e7c548a7b6f",
                "name": "example.com",
                "type": "A",
                "zone_id": "mk2b6fa491c12445a4376666a32429e1",
                "zone_name": "example.com"
            }
        ],
        "success": true
    }`)

	var resp DNSRecordResponse
	err := json.NewDecoder(s).Decode(&resp)
	if err != nil {
		t.Error(err.Error())
	}
	if resp.Success != true {
		t.Errorf("Success Error: %#v != true ", resp.Success)
	}
	if resp.Records[0].ID != "F11cc63e02a42d38174b8e7c548a7b6f" {
		t.Errorf("ID Error: %#v != F11cc63e02a42d38174b8e7c548a7b6f ", resp.Records[0].ID)
	}
	if resp.Records[0].Name != "example.com" {
		t.Errorf("Name Error: %#v != example.com", resp.Records[0].Name)
	}
}
func TestDNSUpdateResponseToJSON(t *testing.T) {
	s := strings.NewReader(`
    {
        "result": {
            "id": "F11cc63e02a42d38174b8e7c548a7b6f",
            "type": "A",
            "name": "example.com",
            "content": "127.0.0.1",
            "proxiable": true,
            "proxied": true,
            "ttl": 1,
            "locked": false,
            "zone_id": "mk2b6fa491c12445a4376666a32429e1",
            "zone_name": "example.com",
            "modified_on": "2018-10-12T14:29:53.205191Z",
            "created_on": "2018-10-12T14:29:53.205191Z",
            "meta": {
              "auto_added": false,
              "managed_by_apps": false,
              "managed_by_argo_tunnel": false
            }
        },
        "success": true,
        "errors": [],
        "messages": []
    }`)

	var resp DNSRecordUpdateResponse
	err := json.NewDecoder(s).Decode(&resp)
	if err != nil {
		t.Error(err.Error())
	}
	if resp.Success != true {
		t.Errorf("Success Error: %#v != true ", resp.Success)
	}
	if resp.Record.ID != "F11cc63e02a42d38174b8e7c548a7b6f" {
		t.Errorf("ID Error: %#v != F11cc63e02a42d38174b8e7c548a7b6f ", resp.Record.ID)
	}
	if resp.Record.Name != "example.com" {
		t.Errorf("Name Error: %#v != example.com", resp.Record.Name)
	}
}

func TestRecordTracked(t *testing.T) {
	s := strings.NewReader(`
    {
        "errors": [],
        "messages": [],
        "result": [
            {
                "content": "127.0.0.1",
                "id": "F11cc63e02a42d38174b8e7c548a7b6f",
                "name": "example.com",
                "type": "A",
                "zone_id": "mk2b6fa491c12445a4376666a32429e1",
                "zone_name": "example.com"
            },
            {
                "content": "127.0.0.1",
                "id": "G00cc63e02a42d38174b8e7c548a7b6f",
                "name": "www.example.com",
                "type": "A",
                "zone_id": "mk2b6fa491c12445a4376666a32429e1",
                "zone_name": "www.example.com"
            }
        ],
        "success": true
    }`)

	var resp DNSRecordResponse
	err := json.NewDecoder(s).Decode(&resp)
	if err != nil {
		t.Error(err.Error())
	}

	domain := &godns.Domain{
		DomainName: "example.com",
		SubDomains: []string{"www"},
	}

	for _, rec := range resp.Records {
		if recordTracked(domain, &rec) {
			t.Info("Record founded: %+v\r\n", rec.Name)
		}
	}
}
