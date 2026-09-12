package lib

import (
	"fmt"
	"net"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

var unsupportedBindWarn sync.Once

// interfaceDialer returns a dialer whose outgoing connections leave through
// the named network interface.
//
// Two mechanisms are combined. The dialer's LocalAddr is set to an address of
// the interface that matches the requested protocol family, and, where the
// platform supports it, the socket is bound to the device itself before
// connecting (SO_BINDTODEVICE on Linux, IP_BOUND_IF on macOS). The second
// step is what actually forces the route: on Linux the source address alone
// does not influence route selection, so without it packets carrying the
// interface's address would still leave via the default route.
func interfaceDialer(ifaceName, network string, timeout time.Duration) (*net.Dialer, error) {
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return nil, fmt.Errorf("query interface %q: %w", ifaceName, err)
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return nil, fmt.Errorf("addresses of query interface %q: %w", ifaceName, err)
	}

	local := pickLocalAddr(addrs, network)
	if local == nil {
		return nil, fmt.Errorf("query interface %q has no usable %s address", ifaceName, network)
	}

	if !deviceBindSupported {
		unsupportedBindWarn.Do(func() {
			log.Warnf("binding sockets to a device is not supported on this platform; "+
				"query_interface %q only sets the source address, routing must send it out the right link", ifaceName)
		})
	}

	log.Debugf("Binding IP query to interface %s with local address %s", ifaceName, local)

	return &net.Dialer{
		Timeout:   timeout,
		KeepAlive: 30 * time.Second,
		LocalAddr: &net.TCPAddr{IP: local},
		Control:   bindToInterface(iface),
	}, nil
}

// pickLocalAddr chooses a global-unicast address of the requested family
// ("tcp4", "tcp6" or the unconstrained "tcp") from an interface's address
// list. Private (RFC 1918 / ULA) addresses are accepted: an interface behind
// a 1:1 NAT is a legitimate query interface. Link-local and loopback
// addresses are skipped since they cannot reach the internet and, for IPv6,
// cannot be bound without a zone.
func pickLocalAddr(addrs []net.Addr, network string) net.IP {
	wantV4 := network == "tcp4"
	wantV6 := network == "tcp6"

	var fallback net.IP
	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}
		if ip == nil || !ip.IsGlobalUnicast() {
			continue
		}

		isV4 := ip.To4() != nil
		switch {
		case wantV4 && isV4, wantV6 && !isV4:
			return ip
		case !wantV4 && !wantV6:
			// Unconstrained: prefer IPv4, accept IPv6.
			if isV4 {
				return ip
			}
			if fallback == nil {
				fallback = ip
			}
		}
	}
	return fallback
}
