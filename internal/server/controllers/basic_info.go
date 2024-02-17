package controllers

import (
	"strings"

	"github.com/TimothyYe/godns/internal/settings"
	"github.com/TimothyYe/godns/internal/utils"
	"github.com/TimothyYe/godns/pkg/lib"
	"github.com/gofiber/fiber/v2"
)

type BasicInfo struct {
	Version      string            `json:"version"`
	StartTime    int64             `json:"start_time"`
	DomainNum    int               `json:"domain_num"`
	SubDomainNum int               `json:"sub_domain_num"`
	Domains      []settings.Domain `json:"domains"`
	PublicIP     string            `json:"public_ip"`
	IPMode       string            `json:"ip_mode"`
	Provider     string            `json:"provider"`
}

func (c *Controller) GetBasicInfo(ctx *fiber.Ctx) error {
	return ctx.JSON(BasicInfo{
		Version:      utils.Version,
		StartTime:    utils.StartTime,
		DomainNum:    c.getDomains(),
		SubDomainNum: c.GetSubDomains(),
		Domains:      c.config.Domains,
		PublicIP:     lib.GetIPHelperInstance(c.config).GetCurrentIP(),
		IPMode:       strings.ToUpper(c.config.IPType),
		Provider:     c.config.Provider,
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
