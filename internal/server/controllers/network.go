package controllers

import (
	"github.com/TimothyYe/godns/internal/settings"
	"github.com/TimothyYe/godns/internal/utils"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

type NetworkSettings struct {
	IPMode        string           `json:"ip_mode"`
	IPUrls        []string         `json:"ip_urls"`
	IPV6Urls      []string         `json:"ipv6_urls"`
	UseProxy      bool             `json:"use_proxy"`
	SkipSSLVerify bool             `json:"skip_ssl_verify"`
	Socks5Proxy   string           `json:"socks5_proxy"`
	Webhook       settings.Webhook `json:"webhook,omitempty"`
	Resolver      string           `json:"resolver"`
	IPInterface   string           `json:"ip_interface"`
}

func (c *Controller) GetNetworkSettings(ctx *fiber.Ctx) error {
	settings := NetworkSettings{
		IPMode:        c.config.IPType,
		IPUrls:        c.config.IPUrls,
		IPV6Urls:      c.config.IPV6Urls,
		UseProxy:      c.config.UseProxy,
		SkipSSLVerify: c.config.SkipSSLVerify,
		Socks5Proxy:   c.config.Socks5Proxy,
		Webhook:       c.config.Webhook,
		Resolver:      c.config.Resolver,
		IPInterface:   c.config.IPInterface,
	}

	return ctx.JSON(settings)
}

func (c *Controller) UpdateNetworkSettings(ctx *fiber.Ctx) error {
	var settings NetworkSettings
	if err := ctx.BodyParser(&settings); err != nil {
		log.Errorf("Failed to parse request body: %s", err.Error())
		return ctx.Status(400).SendString(err.Error())
	}

	if settings.IPMode == utils.IPV4 && len(settings.IPUrls) == 0 {
		return ctx.Status(400).SendString("IP URLs cannot be empty")
	}

	if settings.IPMode == utils.IPV6 && len(settings.IPV6Urls) == 0 {
		return ctx.Status(400).SendString("IPv6 URLs cannot be empty")
	}

	c.config.IPType = settings.IPMode
	if settings.IPMode == utils.IPV6 {
		c.config.IPV6Urls = settings.IPV6Urls
	} else {
		c.config.IPUrls = settings.IPUrls
	}

	c.config.UseProxy = settings.UseProxy
	c.config.SkipSSLVerify = settings.SkipSSLVerify
	c.config.Socks5Proxy = settings.Socks5Proxy
	c.config.Webhook = settings.Webhook
	c.config.Resolver = settings.Resolver
	c.config.IPInterface = settings.IPInterface

	if err := c.config.SaveSettings(c.configPath); err != nil {
		log.Errorf("Failed to save settings: %s", err.Error())
		return ctx.Status(500).SendString("Failed to save network settings")
	}

	return ctx.JSON(settings)
}
