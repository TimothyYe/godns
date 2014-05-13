package main

import (
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func get_currentIP(url string) (string, error) {
	response, err := http.Get(url)
	defer response.Body.Close()

	if err != nil {
		fmt.Println("Cannot get IP...")
		return "", err
	}

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
		fmt.Println("Failed to get domain list...")
		return -1
	}

	sjson, parse_err := simplejson.NewJson([]byte(response))

	if parse_err != nil {
		fmt.Println(parse_err.Error())
		return -1
	}

	if sjson.Get("status").Get("code").MustString() == "1" {
		domains, _ := sjson.Get("domains").Array()

		fmt.Println(domains)

		for _, d := range domains {
			m := d.(map[string]interface{})
			if m["name"] == name {
				id := m["id"]

				switch t := id.(type) {
				case json.Number:
					ret, _ = t.Int64()
				}
			}
		}
	}

	fmt.Printf("Domain id is: %d", ret)
	return ret
}

func get_subdomain(name string) int64 {

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
		fmt.Println("Post failed...")
		fmt.Println(err.Error())
		return "", err
	}

	resp, _ := ioutil.ReadAll(response.Body)

	fmt.Println(string(resp))
	return string(resp), nil
}
