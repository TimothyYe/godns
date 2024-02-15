package controllers

import (
	"github.com/TimothyYe/godns/internal/settings"
	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	config     *settings.Settings
	configPath string
}

func NewController(conf *settings.Settings, configPath string) *Controller {
	return &Controller{
		config:     conf,
		configPath: configPath,
	}
}

func (c *Controller) Auth(ctx *fiber.Ctx) error {
	return ctx.SendString("OK")
}
