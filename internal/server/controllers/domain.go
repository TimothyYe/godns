package controllers

import (
	"github.com/TimothyYe/godns/internal/settings"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

func (c *Controller) GetDomains(ctx *fiber.Ctx) error {
	return ctx.JSON(c.config.Domains)
}

func (c *Controller) AddDomain(ctx *fiber.Ctx) error {
	domain := settings.Domain{}
	if err := ctx.BodyParser(&domain); err != nil {
		log.Errorf("Failed to parse request body: %s", err.Error())
		return ctx.Status(400).SendString(err.Error())
	}

	c.config.Domains = append(c.config.Domains, domain)
	if err := c.config.SaveSettings(c.configPath); err != nil {
		log.Errorf("Failed to save settings: %s", err.Error())
		return ctx.Status(500).SendString("Failed to save settings")
	}

	return ctx.JSON(c.config.Domains)
}
