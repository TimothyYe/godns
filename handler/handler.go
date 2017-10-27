package handler

import "github.com/TimothyYe/godns"

// IHandler is the interface for all DNS handlers
type IHandler interface {
	SetConfiguration(*godns.Settings)
	DomainLoop(domain *godns.Domain, panicChan chan<- godns.Domain)
}

// CreateHandler creates dns handler by different providers
func CreateHandler(provider string) IHandler {
	var handler IHandler

	switch provider {
	case godns.DNSPOD:
		handler = IHandler(&DNSPodHandler{})
	case godns.HE:
		handler = IHandler(&HEHandler{})
	}

	return handler
}
