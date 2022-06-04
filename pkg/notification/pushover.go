package notification

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/TimothyYe/godns/internal/settings"
	"github.com/TimothyYe/godns/internal/utils"

	log "github.com/sirupsen/logrus"
)

const ReqURL = "https://api.pushover.net/1/messages.json"

type PushoverNotification struct {
	conf *settings.Settings
}

func NewPushoverNotification(conf *settings.Settings) INotification {
	return &PushoverNotification{conf: conf}
}

func (n *PushoverNotification) Send(domain, currentIP string) error {
	if n.conf.Notify.Pushover.Token == "" {
		return errors.New("pushover api token cannot be empty")
	}

	if n.conf.Notify.Pushover.User == "" {
		return errors.New("pushover user cannot be empty")
	}

	client := utils.GetHTTPClient(n.conf, false)
	tpl := n.conf.Notify.Pushover.MsgTemplate
	if tpl == "" {
		tpl = "Your IP address changed to <b>{{ .CurrentIP }}</b>. The DNS record for {{ .Domain }} is updated."
		n.conf.Notify.Pushover.HTML = 1
	}

	msg := buildTemplate(currentIP, domain, tpl)
	var response *http.Response
	var err error

	form := url.Values{}
	form.Add("token", n.conf.Notify.Pushover.Token)
	form.Add("user", n.conf.Notify.Pushover.User)
	form.Add("message", msg)
	form.Add("html", strconv.FormatInt(int64(n.conf.Notify.Pushover.HTML), 10))
	if n.conf.Notify.Pushover.Device != "" {
		form.Add("device", n.conf.Notify.Pushover.Device)
	}
	if n.conf.Notify.Pushover.Title != "" {
		form.Add("title", n.conf.Notify.Pushover.Title)
	}
	priority := n.conf.Notify.Pushover.Priority
	if priority != 0 {
		form.Add("priority", strconv.FormatInt(int64(priority), 10))
	}

	log.Debugf("Pushover api request URL: %s, Form: %v", ReqURL, form)
	response, err = client.PostForm(ReqURL, form)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	type APIResponse struct {
		Status  int      `json:"status"`
		Request string   `json:"request"`
		Errors  []string `json:"errors"`
		User    string   `json:"user"`
		Token   string   `json:"token"`
	}
	var resp APIResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		fmt.Println("error:", err)
		return errors.New("failed to parse pushover api response")
	}
	log.Debugf("Pushover api response: %+v", resp)
	if resp.Status != 1 {
		return fmt.Errorf("pushover api call failed Status: %v, Errors: %v", resp.Status, resp.Errors)
	}

	return nil
}
