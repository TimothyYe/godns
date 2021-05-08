package godns

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Domain struct
type Domain struct {
	DomainName string   `json:"domain_name"`
	SubDomains []string `json:"sub_domains"`
}

// Notify struct for slack notification
type SlackNotify struct {
	Enabled     bool   `json:"enabled"`
	BotApiToken string `json:"bot_api_token"`
	Channel     string `json:"channel"`
	MsgTemplate string `json:"message_template"`
	UseProxy    bool   `json:"use_proxy"`
}

// Notify struct for telegram notification
type TelegramNotify struct {
	Enabled     bool   `json:"enabled"`
	BotApiKey   string `json:"bot_api_key"`
	ChatId      string `json:"chat_id"`
	MsgTemplate string `json:"message_template"`
	UseProxy    bool   `json:"use_proxy"`
}

// Notify struct for SMTP notification
type MailNotify struct {
	Enabled      bool   `json:"enabled"`
	SMTPServer   string `json:"smtp_server"`
	SMTPUsername string `json:"smtp_username"`
	SMTPPassword string `json:"smtp_password"`
	SMTPPort     int    `json:"smtp_port"`
	SendTo       string `json:"send_to"`
}

type DiscordNotify struct {
	Enabled     bool   `json:"enabled"`
	BotApiToken string `json:"bot_api_token"`
	Channel     string `json:"channel"`
	MsgTemplate string `json:"message_template"`
}

// Notify struct
type Notify struct {
	Telegram TelegramNotify `json:"telegram"`
	Mail     MailNotify     `json:"mail"`
	Slack    SlackNotify    `json:"slack"`
	Discord  DiscordNotify  `json:"discord"`
}

// Settings struct
type Settings struct {
	Provider    string   `json:"provider"`
	Email       string   `json:"email"`
	Password    string   `json:"password"`
	LoginToken  string   `json:"login_token"`
	Domains     []Domain `json:"domains"`
	IPUrl       string   `json:"ip_url"`
	IPV6Url     string   `json:"ipv6_url"`
	Interval    int      `json:"interval"`
	UserAgent   string   `json:"user_agent,omitempty"`
	Socks5Proxy string   `json:"socks5_proxy"`
	Notify      Notify   `json:"notify"`
	IPInterface string   `json:"ip_interface"`
	IPType      string   `json:"ip_type"`
	Resolver    string   `json:"resolver"`
	UseProxy    bool     `json:"use_proxy"`
	DebugInfo   bool     `json:"debug_info"`
}

// LoadSettings -- Load settings from config file
func LoadSettings(configPath string, settings *Settings) error {
	// LoadSettings from config file
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Println("Error occurs while reading config file, please make sure config file exists!")
		return err
	}

	err = json.Unmarshal(file, settings)
	if err != nil {
		fmt.Println("Error occurs while unmarshal config file, please make sure config file correct!")
		return err
	}

	if settings.Interval == 0 {
		// set default interval as 5 minutes if interval is 0
		settings.Interval = 5 * 60
	}

	return nil
}
