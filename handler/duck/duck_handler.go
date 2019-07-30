package duck

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/TimothyYe/godns"
	"golang.org/x/net/proxy"
)

var (
	// DuckUrl the API address for Duck DNS
	DuckUrl = "https://www.duckdns.org/update?domains=%s&token=%s&ip=%s"
)

// Handler struct
type Handler struct {
	Configuration *godns.Settings
}

// SetConfiguration pass dns settings and store it to handler instance
func (handler *Handler) SetConfiguration(conf *godns.Settings) {
	handler.Configuration = conf
}

// DomainLoop the main logic loop
func (handler *Handler) DomainLoop(domain *godns.Domain, panicChan chan<- godns.Domain) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Recovered in %v: %v\n", err, debug.Stack())
			panicChan <- *domain
		}
	}()

	var lastIP string

	for {
		currentIP, err := godns.GetCurrentIP(handler.Configuration)

		if err != nil {
			log.Println("get_currentIP:", err)
			continue
		}
		log.Println("currentIP is:", currentIP)

		//check against locally cached IP, if no change, skip update
		if currentIP == lastIP {
			log.Printf("IP is the same as cached one. Skip update.\n")
		} else {
			lastIP = currentIP
			client := &http.Client{}

			if handler.Configuration.Socks5Proxy != "" {
				log.Println("use socks5 proxy:" + handler.Configuration.Socks5Proxy)
				dialer, err := proxy.SOCKS5("tcp", handler.Configuration.Socks5Proxy, nil, proxy.Direct)
				if err != nil {
					log.Println("can't connect to the proxy:", err)
					return
				}

				httpTransport := &http.Transport{}
				client.Transport = httpTransport
				httpTransport.Dial = dialer.Dial
			}

			for _, subDomain := range domain.SubDomains {
				// update IP with HTTP GET request
				resp, err := client.Get(fmt.Sprintf(DuckUrl, subDomain, handler.Configuration.LoginToken, currentIP))
				if err != nil {
					// handle error
					log.Print("Failed to update sub domain:", subDomain)
					continue
				}

				defer resp.Body.Close()

				body, err := ioutil.ReadAll(resp.Body)
				if err != nil || string(body) != "OK" {
					// handle error
					log.Print("Failed to update sub domain:", subDomain, err.Error())
					continue
				} else {
					log.Print("IP updated to:", currentIP)
				}

				// Send mail notification if notify is enabled
				if handler.Configuration.Notify.Enabled {
					log.Print("Sending notification to:", handler.Configuration.Notify.SendTo)
					if err := godns.SendNotify(handler.Configuration, fmt.Sprintf("%s.%s", subDomain, domain.DomainName), currentIP); err != nil {
						log.Println("Failed to send notification")
					}
				}
			}
		}

		// Sleep with interval
		log.Printf("Going to sleep, will start next checking in %d seconds...\r\n", handler.Configuration.Interval)
		time.Sleep(time.Second * time.Duration(handler.Configuration.Interval))
	}
}
