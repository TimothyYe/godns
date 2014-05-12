package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func GetCurrentIP(url string) (string, error) {
	response, err := http.Get(url)
	defer response.Body.Close()

	if err != nil {
		fmt.Println("Cannot get IP...")
		return "", err
	}

	body, _ := ioutil.ReadAll(response.Body)
	return string(body), nil
}

func generate_header(content url.Values, setting Settings) url.Values {
	header := url.Values{}
	header.Add("login_email", setting.Email)
	header.Add("login_password", setting.Password)
	header.Add("format", "json")
	header.Add("lang", "en")
	header.Add("error_on_empty", "no")

	for k, _ := range content {
		header.Add(k, content.Get(k))
	}

	return header
}

func post_data(url string, content url.Values, setting Settings) (string, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "https://dnsapi.cn"+url, strings.NewReader(content.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", fmt.Sprintf("GoDNS/0.1 (%s)", setting.Email))

	response, err := client.Do(req)
	defer response.Body.Close()

	if err != nil {
		fmt.Println("Post failed...")
		return "", err
	}

	resp, _ := ioutil.ReadAll(response.Body)

	return string(resp), nil
}
