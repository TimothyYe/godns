package alidns

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// AliDNS token
type AliDNS struct {
	AccessKeyID     string
	AccessKeySecret string
}

var (
	publicParm = map[string]string{
		"AccessKeyId":      "",
		"Format":           "JSON",
		"Version":          "2015-01-09",
		"SignatureMethod":  "HMAC-SHA1",
		"Timestamp":        "",
		"SignatureVersion": "1.0",
		"SignatureNonce":   "",
	}
	baseURL  = "http://alidns.aliyuncs.com/"
	instance *AliDNS
	once     sync.Once
)

type domainRecordsResp struct {
	RequestID     string `json:"RequestId"`
	TotalCount    int
	PageNumber    int
	PageSize      int
	DomainRecords domainRecords
}

type domainRecords struct {
	Record []DomainRecord
}

// DomainRecord struct
type DomainRecord struct {
	DomainName string
	RecordID   string `json:"RecordId"`
	RR         string
	Type       string
	Value      string
	Line       string
	Priority   int
	TTL        int
	Status     string
	Locked     bool
}

func getHTTPBody(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusOK {
		return body, err
	}
	return nil, fmt.Errorf("Status %d, Error:%s", resp.StatusCode, body)
}

func NewAliDNS(key, secret string) *AliDNS {
	once.Do(func() {
		instance = &AliDNS{
			AccessKeyID:     key,
			AccessKeySecret: secret,
		}
	})
	return instance
}

// GetDomainRecords gets all the doamin records according to input subdomain key
func (d *AliDNS) GetDomainRecords(domain, rr string) []DomainRecord {
	resp := &domainRecordsResp{}
	parms := map[string]string{
		"Action":     "DescribeDomainRecords",
		"DomainName": domain,
		"RRKeyWord":  rr,
	}
	urlPath := d.genRequestURL(parms)
	body, err := getHTTPBody(urlPath)
	if err != nil {
		fmt.Printf("GetDomainRecords error.%+v\n", err)
	} else {
		json.Unmarshal(body, resp)
		return resp.DomainRecords.Record
	}
	return nil
}

// UpdateDomainRecord updates domain record
func (d *AliDNS) UpdateDomainRecord(r DomainRecord) error {
	parms := map[string]string{
		"Action":   "UpdateDomainRecord",
		"RecordId": r.RecordID,
		"RR":       r.RR,
		"Type":     r.Type,
		"Value":    r.Value,
		"TTL":      strconv.Itoa(r.TTL),
		"Line":     r.Line,
	}

	urlPath := d.genRequestURL(parms)
	_, err := getHTTPBody(urlPath)
	if err != nil {
		fmt.Printf("UpdateDomainRecord error.%+v\n", err)
	}
	return err
}

func (d *AliDNS) genRequestURL(parms map[string]string) string {
	pArr := []string{}
	ps := map[string]string{}
	for k, v := range publicParm {
		ps[k] = v
	}
	for k, v := range parms {
		ps[k] = v
	}
	now := time.Now().UTC()
	ps["AccessKeyId"] = d.AccessKeyID
	ps["SignatureNonce"] = strconv.Itoa(int(now.UnixNano()) + rand.Intn(99999))
	ps["Timestamp"] = now.Format("2006-01-02T15:04:05Z")

	for k, v := range ps {
		pArr = append(pArr, fmt.Sprintf("%s=%s", k, v))
	}
	sort.Strings(pArr)
	path := strings.Join(pArr, "&")

	s := "GET&%2F&" + url.QueryEscape(path)
	s = strings.Replace(s, "%3A", "%253A", -1)
	s = strings.Replace(s, "%40", "%2540", -1)
	s = strings.Replace(s, "%2A", "%252A", -1)
	mac := hmac.New(sha1.New, []byte(d.AccessKeySecret+"&"))

	mac.Write([]byte(s))
	sign := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return fmt.Sprintf("%s?%s&Signature=%s", baseURL, path, url.QueryEscape(sign))
}
