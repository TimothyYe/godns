package settings

import (
	"errors"
	"fmt"
	"os"
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

// Domain struct.
type Domain struct {
	DomainName string   `json:"domain_name" yaml:"domain_name"`
	SubDomains []string `json:"sub_domains" yaml:"sub_domains"`
	Provider   string   `json:"provider,omitempty" yaml:"provider,omitempty"`
}

// SlackNotify struct for Slack notification.
type SlackNotify struct {
	Enabled         bool   `json:"enabled" yaml:"enabled"`
	BotAPIToken     string `json:"bot_api_token" yaml:"bot_api_token"`
	BotAPITokenFile string `json:"bot_api_token_file" yaml:"bot_api_token_file"`
	Channel         string `json:"channel" yaml:"channel"`
	MsgTemplate     string `json:"message_template" yaml:"message_template"`
}

// TelegramNotify struct for telegram notification.
type TelegramNotify struct {
	Enabled       bool   `json:"enabled" yaml:"enabled"`
	BotAPIKey     string `json:"bot_api_key" yaml:"bot_api_key"`
	BotAPIKeyFile string `json:"bot_api_key_file" yaml:"bot_api_key_file"`
	ChatID        string `json:"chat_id" yaml:"chat_id"`
	MsgTemplate   string `json:"message_template" yaml:"message_template"`
}

// MailNotify struct for SMTP notification.
type MailNotify struct {
	Enabled          bool   `json:"enabled" yaml:"enabled"`
	SMTPServer       string `json:"smtp_server" yaml:"smtp_server"`
	SMTPUsername     string `json:"smtp_username" yaml:"smtp_username"`
	SMTPPassword     string `json:"smtp_password" yaml:"smtp_password"`
	SMTPPasswordFile string `json:"smtp_password_file" yaml:"smtp_password_file"`
	SMTPPort         int    `json:"smtp_port" yaml:"smtp_port"`
	SendFrom         string `json:"send_from" yaml:"send_from"`
	SendTo           string `json:"send_to" yaml:"send_to"`
}

type DiscordNotify struct {
	Enabled         bool   `json:"enabled" yaml:"enabled"`
	BotAPIToken     string `json:"bot_api_token" yaml:"bot_api_token"`
	BotAPITokenFile string `json:"bot_api_token_file" yaml:"bot_api_token_file"`
	Channel         string `json:"channel" yaml:"channel"`
	MsgTemplate     string `json:"message_template" yaml:"message_template"`
}

type PushoverNotify struct {
	Enabled     bool   `json:"enabled" yaml:"enabled"`
	Token       string `json:"token" yaml:"token"`
	TokenFile   string `json:"token_file" yaml:"token:_file"`
	User        string `json:"user" yaml:"user"`
	MsgTemplate string `json:"message_template" yaml:"message_template"`
	Device      string `json:"device" yaml:"device"`
	Title       string `json:"title" yaml:"title"`
	Priority    int    `json:"priority" yaml:"priority"`
	HTML        int    `json:"html" yaml:"html"`
}

type BarkNotify struct {
	Enabled    bool   `json:"enabled" yaml:"enabled"`
	Server     string `json:"server" yaml:"server"`
	Title      string `json:"title" yaml:"title"`
	Subtitle   string `json:"subtitle" yaml:"subtitle"`
	Body       string `json:"body" yaml:"body"`
	DeviceKeys string `json:"device_keys" yaml:"device_keys"`
	Params     string `json:"params" yaml:"params"`
	User       string `json:"user" yaml:"user"`
	Password   string `json:"password" yaml:"password"`
}

// NtfyNotify struct for ntfy notification.
type NtfyNotify struct {
	Enabled     bool   `json:"enabled" yaml:"enabled"`
	Topic       string `json:"topic" yaml:"topic"`
	Server      string `json:"server" yaml:"server"`
	Token       string `json:"token" yaml:"token"`
	User        string `json:"user" yaml:"user"`
	Password    string `json:"password" yaml:"password"`
	Priority    string `json:"priority" yaml:"priority"`
	Tags        string `json:"tags" yaml:"tags"`
	Icon        string `json:"icon" yaml:"icon"`
	MsgTemplate string `json:"message_template" yaml:"message_template"`
}

// ProviderConfig holds provider-specific configuration.
type ProviderConfig struct {
	// Common fields across providers
	Email          string `json:"email,omitempty" yaml:"email,omitempty"`
	Password       string `json:"password,omitempty" yaml:"password,omitempty"`
	PasswordFile   string `json:"password_file,omitempty" yaml:"password_file,omitempty"`
	LoginToken     string `json:"login_token,omitempty" yaml:"login_token,omitempty"`
	LoginTokenFile string `json:"login_token_file,omitempty" yaml:"login_token_file,omitempty"`

	// Provider-specific fields
	AppKey      string `json:"app_key,omitempty" yaml:"app_key,omitempty"`
	AppSecret   string `json:"app_secret,omitempty" yaml:"app_secret,omitempty"`
	ConsumerKey string `json:"consumer_key,omitempty" yaml:"consumer_key,omitempty"`
}

// Notify struct.
type Notify struct {
	Telegram TelegramNotify `json:"telegram" yaml:"telegram"`
	Mail     MailNotify     `json:"mail" yaml:"mail"`
	Slack    SlackNotify    `json:"slack" yaml:"slack"`
	Discord  DiscordNotify  `json:"discord" yaml:"discord"`
	Pushover PushoverNotify `json:"pushover" yaml:"pushover"`
	Bark     BarkNotify     `json:"bark" yaml:"bark"`
	Ntfy     NtfyNotify     `json:"ntfy" yaml:"ntfy"`
}

// Webhook struct.
type Webhook struct {
	Enabled     bool   `json:"enabled" yaml:"enabled"`
	URL         string `json:"url" yaml:"url"`
	RequestBody string `json:"request_body" yaml:"request_body"`
}

type WebPanel struct {
	Enabled  bool   `json:"enabled" yaml:"enabled"`
	Addr     string `json:"addr" yaml:"addr"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

type Mikrotik struct {
	Enabled   bool   `json:"enabled" yaml:"enabled"`
	Addr      string `json:"addr" yaml:"addr"`
	Username  string `json:"username" yaml:"username"`
	Password  string `json:"password" yaml:"password"`
	Interface string `json:"interface" yaml:"interface"`
}

// Settings struct.
type Settings struct {
	// Legacy single provider fields (for backward compatibility)
	Provider       string `json:"provider" yaml:"provider"`
	Email          string `json:"email" yaml:"email"`
	Password       string `json:"password" yaml:"password"`
	PasswordFile   string `json:"password_file" yaml:"password_file"`
	LoginToken     string `json:"login_token" yaml:"login_token"`
	LoginTokenFile string `json:"login_token_file" yaml:"login_token_file"`
	AppKey         string `json:"app_key" yaml:"app_key"`
	AppSecret      string `json:"app_secret" yaml:"app_secret"`
	ConsumerKey    string `json:"consumer_key" yaml:"consumer_key"`

	// New multi-provider configuration
	Providers map[string]*ProviderConfig `json:"providers,omitempty" yaml:"providers,omitempty"`

	// Domain configuration
	Domains []Domain `json:"domains" yaml:"domains"`

	// Network and IP configuration
	IPUrl       string   `json:"ip_url" yaml:"ip_url"`
	IPUrls      []string `json:"ip_urls" yaml:"ip_urls"`
	IPV6Url     string   `json:"ipv6_url" yaml:"ipv6_url"`
	IPV6Urls    []string `json:"ipv6_urls" yaml:"ipv6_urls"`
	IPInterface string   `json:"ip_interface" yaml:"ip_interface"`
	IPType      string   `json:"ip_type" yaml:"ip_type"`
	Resolver    string   `json:"resolver" yaml:"resolver"`

	// Application configuration
	Interval      int    `json:"interval" yaml:"interval"`
	UserAgent     string `json:"user_agent,omitempty" yaml:"user_agent,omitempty"`
	Socks5Proxy   string `json:"socks5_proxy" yaml:"socks5_proxy"`
	UseProxy      bool   `json:"use_proxy" yaml:"use_proxy"`
	DebugInfo     bool   `json:"debug_info" yaml:"debug_info"`
	RunOnce       bool   `json:"run_once" yaml:"run_once"`
	Proxied       bool   `json:"proxied" yaml:"proxied"`
	SkipSSLVerify bool   `json:"skip_ssl_verify" yaml:"skip_ssl_verify"`

	// Feature configuration
	Notify   Notify   `json:"notify" yaml:"notify"`
	Webhook  Webhook  `json:"webhook,omitempty" yaml:"webhook,omitempty"`
	Mikrotik Mikrotik `json:"mikrotik" yaml:"mikrotik"`
	WebPanel WebPanel `json:"web_panel" yaml:"web_panel"`
}

// LoadSettings -- Load settings from config file.
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
	content, err := os.ReadFile(configPath)
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

	if err := loadSecretsFromFile(settings); err != nil {
		return err
	}

	// Load provider-specific secrets
	return loadProviderSecretsFromFile(settings)
}

func (s *Settings) SaveSettings(configPath string) error {
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
	var content []byte
	var err error

	switch fileExt {
	case extJSON:
		content, err = json.MarshalIndent(s, "", "  ")
		if err != nil {
			return err
		}
	case extYML:
		fallthrough
	case extYAML:
		content, err = yaml.Marshal(s)
		if err != nil {
			return err
		}
	default:
		return errors.New("invalid extension for config file:" + fileExt)
	}

	return os.WriteFile(configPath, content, 0644)
}

func loadSecretsFromFile(settings *Settings) error {
	var err error

	if settings.Password, err = readSecretFromFile(
		settings.PasswordFile,
		settings.Password,
	); err != nil {
		return fmt.Errorf("failed to load password from file: %w", err)
	}

	if settings.LoginToken, err = readSecretFromFile(
		settings.LoginTokenFile,
		settings.LoginToken,
	); err != nil {
		return fmt.Errorf("failed to load login token from file: %w", err)
	}

	if settings.Notify.Slack.BotAPIToken, err = readSecretFromFile(
		settings.Notify.Slack.BotAPITokenFile,
		settings.Notify.Slack.BotAPIToken,
	); err != nil {
		return fmt.Errorf("failed to load slack api token from file: %w", err)
	}

	if settings.Notify.Telegram.BotAPIKey, err = readSecretFromFile(
		settings.Notify.Telegram.BotAPIKeyFile,
		settings.Notify.Telegram.BotAPIKey,
	); err != nil {
		return fmt.Errorf("failed to load telegram bot api key from file: %w", err)
	}

	if settings.Notify.Mail.SMTPPassword, err = readSecretFromFile(
		settings.Notify.Mail.SMTPPasswordFile,
		settings.Notify.Mail.SMTPPassword,
	); err != nil {
		return fmt.Errorf("failed to load smtp password from file: %w", err)
	}

	if settings.Notify.Discord.BotAPIToken, err = readSecretFromFile(
		settings.Notify.Discord.BotAPITokenFile,
		settings.Notify.Discord.BotAPIToken,
	); err != nil {
		return fmt.Errorf("failed to load discord bot api token from file: %w", err)
	}

	if settings.Notify.Pushover.Token, err = readSecretFromFile(
		settings.Notify.Pushover.TokenFile,
		settings.Notify.Pushover.Token,
	); err != nil {
		return fmt.Errorf("failed to load pushover token from file: %w", err)
	}

	return nil
}

func readSecretFromFile(source, value string) (string, error) {
	if source == "" {
		return value, nil
	}

	content, err := os.ReadFile(os.ExpandEnv(source))

	if err != nil {
		return value, err
	}

	return strings.TrimSpace(string(content)), nil
}

// GetProviderConfig returns the configuration for a specific provider.
// Falls back to legacy global configuration if provider-specific config is not found.
func (s *Settings) GetProviderConfig(providerName string) *ProviderConfig {
	// If providers map exists and has the specific provider
	if s.Providers != nil {
		if providerConfig, exists := s.Providers[providerName]; exists {
			return providerConfig
		}
	}

	// Fall back to legacy global configuration
	return &ProviderConfig{
		Email:          s.Email,
		Password:       s.Password,
		PasswordFile:   s.PasswordFile,
		LoginToken:     s.LoginToken,
		LoginTokenFile: s.LoginTokenFile,
		AppKey:         s.AppKey,
		AppSecret:      s.AppSecret,
		ConsumerKey:    s.ConsumerKey,
	}
}

// GetDomainProvider returns the provider for a specific domain.
// Falls back to the global provider if domain doesn't specify one.
func (s *Settings) GetDomainProvider(domain *Domain) string {
	if domain.Provider != "" {
		return domain.Provider
	}
	return s.Provider
}

// IsMultiProvider returns true if the configuration uses multiple providers.
func (s *Settings) IsMultiProvider() bool {
	return len(s.Providers) > 0
}

// loadProviderSecretsFromFile loads secrets from files for provider-specific configurations.
func loadProviderSecretsFromFile(settings *Settings) error {
	if settings.Providers == nil {
		return nil
	}

	for providerName, config := range settings.Providers {
		var err error

		if config.Password, err = readSecretFromFile(
			config.PasswordFile,
			config.Password,
		); err != nil {
			return fmt.Errorf("failed to load password from file for provider %s: %w", providerName, err)
		}

		if config.LoginToken, err = readSecretFromFile(
			config.LoginTokenFile,
			config.LoginToken,
		); err != nil {
			return fmt.Errorf("failed to load login token from file for provider %s: %w", providerName, err)
		}
	}

	return nil
}
