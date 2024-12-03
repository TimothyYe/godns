package lib

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"path"
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

func (helper *IPHelper) getIPFromMikrotik() string {
	u, err := url.Parse(helper.configuration.Mikrotik.Addr)
	if err != nil {
		log.Error("fail to parse mikrotik address: ", err)
		return ""
	}
	u.Path = path.Join(u.Path, "/rest/ip/address")
	q := u.Query()
	q.Add("interface", helper.configuration.Mikrotik.Interface)
	q.Add(".proplist", "address")
	u.RawQuery = q.Encode()

	req, _ := http.NewRequest("GET", u.String(), nil)
	auth := fmt.Sprintf("%s:%s", helper.configuration.Mikrotik.Username, helper.configuration.Mikrotik.Password)
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(auth)))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{
		Timeout:   time.Second * utils.DefaultTimeout,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}

	response, err := client.Do(req)
	if err != nil {
		log.Error("request mikrotik address failed:", err)
		return ""
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Error("request code failed: ", response.StatusCode)
		return ""
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Error("read body failed: ", err)
		return ""
	}

	m := []map[string]string{}
	if err := json.Unmarshal(body, &m); err != nil {
		log.Error("unmarshal body failed: ", err)
		return ""
	}
	if len(m) < 1 {
		log.Error("could not get ip from response: ", m)
		return ""
	}

	res := strings.Split(m[0]["address"], "/")
	return res[0]
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
			log.Debugf("Get ip success from network interface by: %s, IP: %s", helper.configuration.IPInterface, ip.String())
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

	if helper.configuration.Mikrotik.Enabled {
		ip = helper.getIPFromMikrotik()
		if ip == "" {
			log.Error("get ip from mikrotik failed. Fallback to get ip from onlinke if possible.")
		} else {
			helper.setCurrentIP(ip)
			return
		}
	}

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
	transport := &http.Transport{
		DialContext: func(ctx context.Context, _, addr string) (net.Conn, error) {
			proto := "tcp"

			if strings.ToUpper(helper.configuration.IPType) == utils.IPV4 {
				// Force the network to "tcp4" to use only IPv4
				proto = "tcp4"
			}

			return (&net.Dialer{
				Timeout:   time.Second * utils.DefaultTimeout,
				KeepAlive: 30 * time.Second,
			}).DialContext(ctx, proto, addr)
		},
	}

	client := &http.Client{
		Timeout:   time.Second * utils.DefaultTimeout,
		Transport: transport,
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

		if isIPv4(onlineIP) {
			if strings.ToUpper(helper.configuration.IPType) != utils.IPV4 {
				log.Warnf("The online IP (%s) from %s is not IPV6, will skip it.", onlineIP, reqURL)
				continue
			}
		} else {
			if strings.ToUpper(helper.configuration.IPType) != utils.IPV6 {
				log.Warnf("The online IP (%s) from %s is not IPV4, will skip it.", onlineIP, reqURL)
				continue
			}
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
