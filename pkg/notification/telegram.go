package notification

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/TimothyYe/godns/internal/settings"
	"github.com/TimothyYe/godns/internal/utils"
)

type TelegramNotification struct {
	conf *settings.Settings
}

func NewTelegramNotification(conf *settings.Settings) INotification {
	return &TelegramNotification{conf: conf}
}

func (n *TelegramNotification) Send(domain, currentIP string) error {
	if n.conf.Notify.Telegram.BotAPIKey == "" {
		return errors.New("bot api key cannot be empty")
	}

	if n.conf.Notify.Telegram.ChatID == "" {
		return errors.New("chat id cannot be empty")
	}

	client := utils.GetHTTPClient(n.conf, n.conf.Notify.Telegram.UseProxy)
	tpl := n.conf.Notify.Telegram.MsgTemplate
	if tpl == "" {
		tpl = "_Your IP address is changed to_%0A%0A*{{ .CurrentIP }}*%0A%0ADomain *{{ .Domain }}* is updated"
	}

	msg := buildTemplate(currentIP, domain, tpl)
	reqURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&parse_mode=Markdown&text=%s",
		n.conf.Notify.Telegram.BotAPIKey,
		n.conf.Notify.Telegram.ChatID,
		msg)
	var response *http.Response
	var err error

	response, err = client.Get(reqURL)

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
