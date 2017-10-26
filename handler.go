package main

//IHandler is the interface for all DNS handlers
type IHandler interface {
	DomainLoop(domain *Domain, panicChan chan<- Domain)
}

func createHandler(provider string) IHandler {
	var handler IHandler

	switch provider {
	case DNSPOD:
		handler = IHandler(&DNSPodHandler{})
	case HE:
		handler = IHandler(&HEHandler{})
	}

	return handler
}
