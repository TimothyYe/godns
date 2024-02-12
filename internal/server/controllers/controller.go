package controllers

import (
	"github.com/TimothyYe/godns/internal/settings"
	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	config *settings.Settings
}

func NewController(conf *settings.Settings) *Controller {
	return &Controller{
		config: conf,
	}
}

func (c *Controller) Auth(ctx *fiber.Ctx) error {
	return ctx.SendString("OK")
}
