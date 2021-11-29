package settings

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"gopkg.in/yaml.v3"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

const (
	extJSON = "json"
	extYAML = "yaml"
	extYML  = "yml"
)

// Domain struct
type Domain struct {
	DomainName string   `json:"domain_name" yaml:"domain_name"`
	SubDomains []string `json:"sub_domains" yaml:"sub_domains"`
}

// SlackNotify struct for Slack notification
type SlackNotify struct {
	Enabled     bool   `json:"enabled" yaml:"enabled"`
	BotApiToken string `json:"bot_api_token" yaml:"bot_api_token"`
	Channel     string `json:"channel" yaml:"channel"`
	MsgTemplate string `json:"message_template" yaml:"message_template"`
	UseProxy    bool   `json:"use_proxy" yaml:"use_proxy"`
}

// TelegramNotify struct for telegram notification
type TelegramNotify struct {
	Enabled     bool   `json:"enabled" yaml:"enabled"`
	BotApiKey   string `json:"bot_api_key" yaml:"bot_api_key"`
	ChatId      string `json:"chat_id" yaml:"chat_id"`
	MsgTemplate string `json:"message_template" yaml:"message_template"`
	UseProxy    bool   `json:"use_proxy" yaml:"use_proxy"`
}

// MailNotify struct for SMTP notification
type MailNotify struct {
	Enabled      bool   `json:"enabled" yaml:"enabled"`
	SMTPServer   string `json:"smtp_server" yaml:"smtp_server"`
	SMTPUsername string `json:"smtp_username" yaml:"smtp_username"`
	SMTPPassword string `json:"smtp_password" yaml:"smtp_password"`
	SMTPPort     int    `json:"smtp_port" yaml:"smtp_port"`
	SendTo       string `json:"send_to" yaml:"send_to"`
}

type DiscordNotify struct {
	Enabled     bool   `json:"enabled" yaml:"enabled"`
	BotApiToken string `json:"bot_api_token" yaml:"bot_api_token"`
	Channel     string `json:"channel" yaml:"channel"`
	MsgTemplate string `json:"message_template" yaml:"message_template"`
}

// Notify struct
type Notify struct {
	Telegram TelegramNotify `json:"telegram" yaml:"telegram"`
	Mail     MailNotify     `json:"mail" yaml:"mail"`
	Slack    SlackNotify    `json:"slack" yaml:"slack"`
	Discord  DiscordNotify  `json:"discord" yaml:"discord"`
}

// Settings struct
type Settings struct {
	Provider    string   `json:"provider" yaml:"provider"`
	Email       string   `json:"email" yaml:"email"`
	Password    string   `json:"password" yaml:"password"`
	LoginToken  string   `json:"login_token" yaml:"login_token"`
	Domains     []Domain `json:"domains" yaml:"domains"`
	IPUrl       string   `json:"ip_url" yaml:"ip_url"`
	IPV6Url     string   `json:"ipv6_url" yaml:"ipv6_url"`
	Interval    int      `json:"interval" yaml:"interval"`
	UserAgent   string   `json:"user_agent,omitempty" yaml:"user_agent,omitempty"`
	Socks5Proxy string   `json:"socks5_proxy" yaml:"socks5_proxy"`
	Notify      Notify   `json:"notify" yaml:"notify"`
	IPInterface string   `json:"ip_interface" yaml:"ip_interface"`
	IPType      string   `json:"ip_type" yaml:"ip_type"`
	Resolver    string   `json:"resolver" yaml:"resolver"`
	UseProxy    bool     `json:"use_proxy" yaml:"use_proxy"`
	DebugInfo   bool     `json:"debug_info" yaml:"debug_info"`
}

// LoadSettings -- Load settings from config file
func LoadSettings(configPath string, settings *Settings) error {
	// get config file extension
	fileExt := strings.ToLower(filepath.Ext(configPath))
	if fileExt == "" {
		return errors.New("invalid file extension")
	}

	// get file name without extension
	fileName := strings.TrimSuffix(filepath.Base(configPath), fileExt)
	fileExt = fileExt[1:]

	if fileName == "" {
		return errors.New("invalid config file name")
	}

	// LoadSettings from config file
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Println("Error occurs while reading config file, please make sure config file exists!")
		return err
	}

	switch fileExt {
	case extJSON:
		if err := json.Unmarshal(content, settings); err != nil {
			return err
		}
	case extYML:
		fallthrough
	case extYAML:
		if err := yaml.Unmarshal(content, settings); err != nil {
			return err
		}
	default:
		return errors.New("invalid extension for config file:" + fileExt)
	}

	if settings.Interval == 0 {
		// set default interval as 5 minutes if interval is 0
		settings.Interval = 5 * 60
	}

	return nil
}
