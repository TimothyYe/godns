package godns

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strings"

	"golang.org/x/net/proxy"
)

const (
	// PanicMax is the max allowed panic times
	PanicMax = 5
	// INTERVAL is minute
	INTERVAL = 5
	// DNSPOD for dnspod.cn
	DNSPOD = "DNSPod"
	// HE for he.net
	HE = "HE"
)

func GetCurrentIP(configuration *Settings) (string, error) {
	client := &http.Client{}

	if configuration.Socks5Proxy != "" {

		log.Println("use socks5 proxy:" + configuration.Socks5Proxy)

		dialer, err := proxy.SOCKS5("tcp", configuration.Socks5Proxy, nil, proxy.Direct)
		if err != nil {
			fmt.Println("can't connect to the proxy:", err)
			return "", err
		}

		httpTransport := &http.Transport{}
		client.Transport = httpTransport
		httpTransport.Dial = dialer.Dial
	}

	response, err := client.Get(configuration.IPUrl)

	if err != nil {
		log.Println("Cannot get IP...")
		return "", err
	}

	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	return string(body), nil
}

func IdentifyPanic() string {
	var name, file string
	var line int
	var pc [16]uintptr

	n := runtime.Callers(3, pc[:])
	for _, pc := range pc[:n] {
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}
		file, line = fn.FileLine(pc)
		name = fn.Name()
		if !strings.HasPrefix(name, "runtime.") {
			break
		}
	}

	switch {
	case name != "":
		return fmt.Sprintf("%v:%v", name, line)
	case file != "":
		return fmt.Sprintf("%v:%v", file, line)
	}

	return fmt.Sprintf("pc:%x", pc)
}

func Usage() {
	log.Println("[command] -c=[config file path]")
	flag.PrintDefaults()
}

func CheckSettings(config *Settings) error {
	if config.Provider == DNSPOD {
		if (config.Email == "" || config.Password == "") && config.LoginToken == "" {
			return errors.New("email/password or login token cannot be empty")
		}
	} else if config.Provider == HE {
		if config.Password == "" {
			return errors.New("password cannot be empty")
		}
	} else {
		return errors.New("please provide supported DNS provider: DNSPod/HE")
	}

	return nil
}
