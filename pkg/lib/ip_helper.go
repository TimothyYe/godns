package lib

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/TimothyYe/godns/internal/utils"

	"github.com/TimothyYe/godns/internal/settings"
)

type IPHelper struct {
	reqURLs       []string
	currentIP     string
	mutex         sync.RWMutex
	configuration *settings.Settings
	idx           int64
}

var (
	helperInstance *IPHelper
	helperOnce     sync.Once
)

func (helper *IPHelper) UpdateConfiguration(conf *settings.Settings) {
	helper.mutex.Lock()
	defer helper.mutex.Unlock()

	// clear urls
	helper.reqURLs = helper.reqURLs[:0]
	// reset the index
	helper.idx = -1

	if conf.IPType == "" || strings.ToUpper(conf.IPType) == utils.IPV4 {
		// filter empty urls
		for _, url := range conf.IPUrls {
			if url != "" {
				helper.reqURLs = append(helper.reqURLs, url)
			}
		}

		if conf.IPUrl != "" {
			helper.reqURLs = append(helper.reqURLs, conf.IPUrl)
		}
	} else {
		// filter empty urls
		for _, url := range conf.IPV6Urls {
			if url != "" {
				helper.reqURLs = append(helper.reqURLs, url)
			}
		}

		if conf.IPV6Url != "" {
			helper.reqURLs = append(helper.reqURLs, conf.IPV6Url)
		}
	}

	log.Debugf("Update ip helper configuration, urls: %v", helper.reqURLs)
}

func GetIPHelperInstance(conf *settings.Settings) *IPHelper {
	helperOnce.Do(func() {
		helperInstance = &IPHelper{
			configuration: conf,
			idx:           -1,
		}

		SafeGo(func() {
			for {
				helperInstance.getCurrentIP()
				time.Sleep(time.Second * time.Duration(conf.Interval))
			}
		})
	})

	return helperInstance
}

func (helper *IPHelper) GetCurrentIP() string {
	// for the first load
	if helper.currentIP == "" {
		helper.getCurrentIP()
	}

	helper.mutex.RLock()
	defer helper.mutex.RUnlock()

	return helper.currentIP
}

func (helper *IPHelper) setCurrentIP(ip string) {
	helper.mutex.Lock()
	defer helper.mutex.Unlock()

	helper.currentIP = ip
}

func (helper *IPHelper) getNext() string {
	newIdx := atomic.AddInt64(&helper.idx, 1)

	helper.mutex.RLock()
	defer helper.mutex.RUnlock()
	newIdx %= int64(len(helper.reqURLs))
	next := helper.reqURLs[newIdx]
	return next
}

// getIPFromInterface gets IP address from the specific interface.
func (helper *IPHelper) getIPFromInterface() (string, error) {
	ifaces, err := net.InterfaceByName(helper.configuration.IPInterface)
	if err != nil {
		log.Error("Can't get network device "+helper.configuration.IPInterface+":", err)
		return "", err
	}

	addrs, err := ifaces.Addrs()
	if err != nil {
		log.Error("Can't get address from "+helper.configuration.IPInterface+":", err)
		return "", err
	}

	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}
		if ip == nil {
			continue
		}

		if ip.IsPrivate() {
			continue
		}

		if isIPv4(ip.String()) {
			if strings.ToUpper(helper.configuration.IPType) != utils.IPV4 {
				continue
			}
		} else {
			if strings.ToUpper(helper.configuration.IPType) != utils.IPV6 {
				continue
			}
		}

		if ip.String() != "" {
			log.Debugf("Get ip success from network intereface by: %s, IP: %s", helper.configuration.IPInterface, ip.String())
			return ip.String(), nil
		}
	}
	return "", errors.New("can't get a valid address from " + helper.configuration.IPInterface)
}

func isIPv4(ip string) bool {
	return strings.Count(ip, ":") < 2
}

// getCurrentIP gets an IP from either internet or specific interface, depending on configuration.
func (helper *IPHelper) getCurrentIP() {
	var err error
	var ip string

	if len(helper.reqURLs) > 0 {
		ip = helper.getIPOnline()
		if ip == "" {
			log.Error("get ip online failed. Fallback to get ip from interface if possible.")
		} else {
			helper.setCurrentIP(ip)
			return
		}
	}

	if helper.configuration.IPInterface != "" {
		ip, err = helper.getIPFromInterface()
		if err != nil {
			log.Error("get ip from interface failed. There is no more ways to try.")
		} else {
			helper.setCurrentIP(ip)
			return
		}
	}
}

// getIPOnline gets public IP from internet.
func (helper *IPHelper) getIPOnline() string {
	client := &http.Client{
		Timeout: time.Second * utils.DefaultTimeout,
	}
	var onlineIP string

	for {
		reqURL := helper.getNext()
		req, _ := http.NewRequest("GET", reqURL, nil)

		if helper.configuration.UserAgent != "" {
			req.Header.Set("User-Agent", helper.configuration.UserAgent)
		}

		response, err := client.Do(req)

		if err != nil {
			log.Error("Cannot get IP:", err)
			time.Sleep(time.Millisecond * 300)
			continue
		}

		if response.StatusCode != http.StatusOK {
			log.Error(fmt.Sprintf("request %v got httpCode:%v", reqURL, response.StatusCode))
			continue
		}

		body, _ := io.ReadAll(response.Body)
		ipReg := regexp.MustCompile(utils.IPPattern)
		onlineIP = ipReg.FindString(string(body))
		if onlineIP == "" {
			log.Error(fmt.Sprintf("request:%v failed to get online IP", reqURL))
			continue
		}
		log.Debugf("Get ip success by: %s, online IP: %s", reqURL, onlineIP)

		err = response.Body.Close()
		if err != nil {
			log.Error(fmt.Sprintf("request:%v failed to get online IP", reqURL))
			continue
		}

		if onlineIP == "" {
			log.Error("fail to get online IP")
		}

		break
	}

	return onlineIP
}
