package controllers

import (
	"github.com/TimothyYe/godns/internal/settings"
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

func (c *Controller) GetMultiProviders(ctx *fiber.Ctx) error {
	// Combine legacy provider config with new multi-provider config
	result := make(map[string]*settings.ProviderConfig)

	// First, add all providers from the new multi-provider map
	if c.config.Providers != nil {
		for name, config := range c.config.Providers {
			result[name] = config
		}
	}

	// Then, add the legacy provider if it exists and isn't already in the map
	if c.config.Provider != "" {
		if _, exists := result[c.config.Provider]; !exists {
			result[c.config.Provider] = &settings.ProviderConfig{
				Email:          c.config.Email,
				Password:       c.config.Password,
				PasswordFile:   c.config.PasswordFile,
				LoginToken:     c.config.LoginToken,
				LoginTokenFile: c.config.LoginTokenFile,
				AppKey:         c.config.AppKey,
				AppSecret:      c.config.AppSecret,
				ConsumerKey:    c.config.ConsumerKey,
			}
		}
	}

	return ctx.JSON(result)
}

func (c *Controller) UpdateMultiProviders(ctx *fiber.Ctx) error {
	var providers map[string]*settings.ProviderConfig
	if err := ctx.BodyParser(&providers); err != nil {
		return err
	}

	// If there's a legacy provider configured, check if it's in the input
	if c.config.Provider != "" {
		if legacyConfig, hasLegacy := providers[c.config.Provider]; hasLegacy {
			// Update the legacy top-level fields
			c.config.Email = legacyConfig.Email
			c.config.Password = legacyConfig.Password
			c.config.PasswordFile = legacyConfig.PasswordFile
			c.config.LoginToken = legacyConfig.LoginToken
			c.config.LoginTokenFile = legacyConfig.LoginTokenFile
			c.config.AppKey = legacyConfig.AppKey
			c.config.AppSecret = legacyConfig.AppSecret
			c.config.ConsumerKey = legacyConfig.ConsumerKey

			// Create a new map with only the non-legacy providers
			nonLegacyProviders := make(map[string]*settings.ProviderConfig)
			for name, config := range providers {
				if name != c.config.Provider {
					nonLegacyProviders[name] = config
				}
			}

			// Update the providers map with only non-legacy providers
			if len(nonLegacyProviders) > 0 {
				c.config.Providers = nonLegacyProviders
			} else {
				// If no other providers, clear the providers map
				c.config.Providers = nil
			}

			if err := c.config.SaveSettings(c.configPath); err != nil {
				log.Errorf("Failed to save settings: %s", err.Error())
				return ctx.Status(500).SendString("Failed to save settings")
			}

			return ctx.SendStatus(fiber.StatusOK)
		}
	}

	// No legacy provider match, update the entire providers map
	c.config.Providers = providers

	if err := c.config.SaveSettings(c.configPath); err != nil {
		log.Errorf("Failed to save settings: %s", err.Error())
		return ctx.Status(500).SendString("Failed to save settings")
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (c *Controller) AddProviderConfig(ctx *fiber.Ctx) error {
	providerName := ctx.Params("provider")
	if providerName == "" {
		return ctx.Status(400).SendString("Provider name is required")
	}

	var config settings.ProviderConfig
	if err := ctx.BodyParser(&config); err != nil {
		return err
	}

	if c.config.Providers == nil {
		c.config.Providers = make(map[string]*settings.ProviderConfig)
	}

	c.config.Providers[providerName] = &config

	if err := c.config.SaveSettings(c.configPath); err != nil {
		log.Errorf("Failed to save settings: %s", err.Error())
		return ctx.Status(500).SendString("Failed to save settings")
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (c *Controller) DeleteProviderConfig(ctx *fiber.Ctx) error {
	providerName := ctx.Params("provider")
	if providerName == "" {
		return ctx.Status(400).SendString("Provider name is required")
	}

	if c.config.Providers != nil {
		delete(c.config.Providers, providerName)

		if err := c.config.SaveSettings(c.configPath); err != nil {
			log.Errorf("Failed to save settings: %s", err.Error())
			return ctx.Status(500).SendString("Failed to save settings")
		}
	}

	return ctx.SendStatus(fiber.StatusOK)
}
