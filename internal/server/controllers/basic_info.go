package controllers

import (
	"strings"

	"github.com/TimothyYe/godns/internal/settings"
	"github.com/TimothyYe/godns/internal/utils"
	"github.com/TimothyYe/godns/pkg/lib"
	"github.com/gofiber/fiber/v2"
)

type BasicInfo struct {
	Version         string            `json:"version"`
	StartTime       int64             `json:"start_time"`
	DomainNum       int               `json:"domain_num"`
	SubDomainNum    int               `json:"sub_domain_num"`
	Domains         []settings.Domain `json:"domains"`
	PublicIP        string            `json:"public_ip"`
	IPMode          string            `json:"ip_mode"`
	Provider        string            `json:"provider"`
	IsMultiProvider bool              `json:"is_multi_provider"`
	Providers       []string          `json:"providers"`
}

func (c *Controller) GetBasicInfo(ctx *fiber.Ctx) error {
	isMultiProvider := c.config.IsMultiProvider()
	providers := c.getProviders()

	return ctx.JSON(BasicInfo{
		Version:         utils.Version,
		StartTime:       utils.StartTime,
		DomainNum:       c.getDomains(),
		SubDomainNum:    c.GetSubDomains(),
		Domains:         c.config.Domains,
		PublicIP:        lib.GetIPHelperInstance(c.config).GetCurrentIP(),
		IPMode:          strings.ToUpper(c.config.IPType),
		Provider:        c.config.Provider,
		IsMultiProvider: isMultiProvider,
		Providers:       providers,
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

func (c *Controller) getProviders() []string {
	providersSet := make(map[string]bool)

	// Add global provider if specified (legacy single provider mode)
	if c.config.Provider != "" {
		providersSet[c.config.Provider] = true
	}

	// Add providers from multi-provider configuration
	if c.config.Providers != nil {
		for providerName := range c.config.Providers {
			providersSet[providerName] = true
		}
	}

	// Add providers from domains (for mixed configuration)
	for _, domain := range c.config.Domains {
		if domain.Provider != "" {
			providersSet[domain.Provider] = true
		}
	}

	// Convert set to slice
	providers := make([]string, 0, len(providersSet))
	for provider := range providersSet {
		providers = append(providers, provider)
	}

	return providers
}
