package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/bitly/go-simplejson"
)

func get_currentIP(url string) (string, error) {
	response, err := http.Get(url)

	if err != nil {
		log.Println("Cannot get IP...")
		return "", err
	}

	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	return string(body), nil
}

func generate_header(content url.Values) url.Values {
	header := url.Values{}
	header.Add("login_email", Configuration.Email)
	header.Add("login_password", Configuration.Password)
	header.Add("format", "json")
	header.Add("lang", "en")
	header.Add("error_on_empty", "no")

	if content != nil {
		for k, _ := range content {
			header.Add(k, content.Get(k))
		}
	}

	return header
}

func api_version() {
	fmt.Println(Configuration.Email)
	post_data("/Info.Version", nil)
}

func get_domain(name string) int64 {

	var ret int64
	values := url.Values{}
	values.Add("type", "all")
	values.Add("offset", "0")
	values.Add("length", "20")

	response, err := post_data("/Domain.List", values)

	if err != nil {
		log.Println("Failed to get domain list...")
		return -1
	}

	sjson, parse_err := simplejson.NewJson([]byte(response))

	if parse_err != nil {
		log.Println(parse_err)
		return -1
	}

	if sjson.Get("status").Get("code").MustString() == "1" {
		domains, _ := sjson.Get("domains").Array()

		for _, d := range domains {
			m := d.(map[string]interface{})
			if m["name"] == name {
				id := m["id"]

				switch t := id.(type) {
				case json.Number:
					ret, _ = t.Int64()
				}

				break
			}
		}
		if len(domains) == 0 {
			log.Println("domains slice is empty.")
		}
	} else {
		log.Println("get_domain:status code:", sjson.Get("status").Get("code").MustString())
	}

	return ret
}

func get_subdomain(domain_id int64, name string) (string, string) {
	log.Println("debug:", domain_id, name)
	var ret, ip string
	value := url.Values{}
	value.Add("domain_id", strconv.FormatInt(domain_id, 10))
	value.Add("offset", "0")
	value.Add("length", "1")
	value.Add("sub_domain", name)

	response, err := post_data("/Record.List", value)

	if err != nil {
		log.Println("Failed to get domain list")
		return "", ""
	}

	sjson, parse_err := simplejson.NewJson([]byte(response))

	if parse_err != nil {
		log.Println(parse_err)
		return "", ""
	}

	if sjson.Get("status").Get("code").MustString() == "1" {
		records, _ := sjson.Get("records").Array()

		for _, d := range records {
			m := d.(map[string]interface{})
			if m["name"] == name {
				ret = m["id"].(string)
				ip = m["value"].(string)
				break
			}
		}
		if len(records) == 0 {
			log.Println("records slice is empty.")
		}
	} else {
		log.Println("get_subdomain:status code:", sjson.Get("status").Get("code").MustString())
	}

	return ret, ip
}

func update_ip(domain_id int64, sub_domain_id string, sub_domain_name string, ip string) {
	value := url.Values{}
	value.Add("domain_id", strconv.FormatInt(domain_id, 10))
	value.Add("record_id", sub_domain_id)
	value.Add("sub_domain", sub_domain_name)
	value.Add("record_type", "A")
	value.Add("record_line", "默认")
	value.Add("value", ip)

	response, err := post_data("/Record.Modify", value)

	if err != nil {
		log.Println("Failed to update record to new IP!")
		log.Println(err)
		return
	}

	sjson, parse_err := simplejson.NewJson([]byte(response))

	if parse_err != nil {
		log.Println(parse_err)
		return
	}

	if sjson.Get("status").Get("code").MustString() == "1" {
		log.Println("New IP updated!")
	}

}

func post_data(url string, content url.Values) (string, error) {
	client := &http.Client{}
	values := generate_header(content)
	req, _ := http.NewRequest("POST", "https://dnsapi.cn"+url, strings.NewReader(values.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", fmt.Sprintf("GoDNS/0.1 (%s)", Configuration.Email))

	response, err := client.Do(req)
	defer response.Body.Close()

	if err != nil {
		log.Println("Post failed...")
		log.Println(err)
		return "", err
	}

	resp, _ := ioutil.ReadAll(response.Body)

	return string(resp), nil
}
