package notify

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/TimothyYe/godns"
)

type SlackNotify struct {
	conf *godns.Settings
}

func NewSlackNotify(conf *godns.Settings) INotify {
	return &SlackNotify{conf: conf}
}

func (n *SlackNotify) Send(domain, currentIP string) error {
	if n.conf.Notify.Slack.BotApiToken == "" {
		return errors.New("bot api token cannot be empty")
	}

	if n.conf.Notify.Slack.Channel == "" {
		return errors.New("channel cannot be empty")
	}
	client := godns.GetHttpClient(n.conf, n.conf.Notify.Slack.UseProxy)
	tpl := n.conf.Notify.Slack.MsgTemplate
	if tpl == "" {
		tpl = "_Your IP address is changed to_\n\n*{{ .CurrentIP }}*\n\nDomain *{{ .Domain }}* is updated"
	}

	msg := buildTemplate(currentIP, domain, tpl)

	var response *http.Response
	var err error

	formData := url.Values{
		"token":   {n.conf.Notify.Slack.BotApiToken},
		"channel": {n.conf.Notify.Slack.Channel},
		"text":    {msg},
	}

	response, err = client.PostForm("https://slack.com/api/chat.postMessage", formData)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	type ResponseParameters struct {
		MigrateToChatID int64 `json:"migrate_to_chat_id"` // optional
		RetryAfter      int   `json:"retry_after"`        // optional
	}
	type APIResponse struct {
		Ok          bool                `json:"ok"`
		Result      json.RawMessage     `json:"result"`
		ErrorCode   int                 `json:"error_code"`
		Description string              `json:"description"`
		Parameters  *ResponseParameters `json:"parameters"`
	}
	var resp APIResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		fmt.Println("error:", err)
		return errors.New("failed to parse response")
	}
	if !resp.Ok {
		return errors.New(resp.Description)
	}

	return nil
}
