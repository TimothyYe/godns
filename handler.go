package main

type IHandler interface {
	DomainLoop(domain *Domain)
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
