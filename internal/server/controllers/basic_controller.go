package controllers

import (
	"github.com/TimothyYe/godns/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type BasicInfo struct {
	Version    string `json:"version"`
	StartTime  int64  `json:"start_time"`
	Domains    int    `json:"domains"`
	SubDomains int    `json:"sub_domains"`
}

func (c *Controller) GetBasicInfo(ctx *fiber.Ctx) error {
	return ctx.JSON(BasicInfo{
		Version:    utils.Version,
		StartTime:  utils.StartTime,
		Domains:    c.getDomains(),
		SubDomains: c.GetSubDomains(),
	})
}

func (c *Controller) getDomains() int {
	// count the total number of domains
	return len(c.config.Domains)
}

func (c *Controller) GetSubDomains() int {
	// get the total number of all the sub domains
	var count int
	for _, domain := range c.config.Domains {
		count += len(domain.SubDomains)
	}

	return count
}
