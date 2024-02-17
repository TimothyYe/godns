package controllers

import (
	"github.com/TimothyYe/godns/internal/utils"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

type Provider struct {
	Provider    string `json:"provider" yaml:"provider"`
	Email       string `json:"email" yaml:"email"`
	Password    string `json:"password" yaml:"password"`
	LoginToken  string `json:"login_token" yaml:"login_token"`
	AppKey      string `json:"app_key" yaml:"app_key"`
	AppSecret   string `json:"app_secret" yaml:"app_secret"`
	ConsumerKey string `json:"consumer_key" yaml:"consumer_key"`
}

func (c *Controller) GetProvider(ctx *fiber.Ctx) error {
	provider := Provider{
		Provider:    c.config.Provider,
		Email:       c.config.Email,
		Password:    c.config.Password,
		LoginToken:  c.config.LoginToken,
		AppKey:      c.config.AppKey,
		AppSecret:   c.config.AppSecret,
		ConsumerKey: c.config.ConsumerKey,
	}
	return ctx.JSON(provider)
}

func (c *Controller) GetProviderSettings(ctx *fiber.Ctx) error {
	return ctx.JSON(utils.Providers)
}

func (c *Controller) UpdateProvider(ctx *fiber.Ctx) error {
	var provider Provider
	if err := ctx.BodyParser(&provider); err != nil {
		return err
	}

	c.config.Provider = provider.Provider
	c.config.Email = provider.Email
	c.config.Password = provider.Password
	c.config.LoginToken = provider.LoginToken
	c.config.AppKey = provider.AppKey
	c.config.AppSecret = provider.AppSecret
	c.config.ConsumerKey = provider.ConsumerKey

	if err := c.config.SaveSettings(c.configPath); err != nil {
		log.Errorf("Failed to save settings: %s", err.Error())
		return ctx.Status(500).SendString("Failed to save settings")
	}

	return ctx.SendStatus(fiber.StatusOK)
}
