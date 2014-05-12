package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
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
