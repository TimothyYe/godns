package controllers

import "github.com/gofiber/fiber/v2"

type Controller struct {
}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) Ping(ctx *fiber.Ctx) error {
	return ctx.SendString("Pong")
}
