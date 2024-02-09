package controllers

import (
	"github.com/TimothyYe/godns/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type BasicInfo struct {
	Version   string   `json:"version"`
	StartTime int64    `json:"start_time"`
	Providers []string `json:"providers"`
}

type Controller struct {
}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) Auth(ctx *fiber.Ctx) error {
	return ctx.SendString("OK")
}

func (c *Controller) GetBasicInfo(ctx *fiber.Ctx) error {
	return ctx.JSON(BasicInfo{
		Version:   utils.Version,
		StartTime: utils.StartTime,
		Providers: utils.Providers,
	})
}
