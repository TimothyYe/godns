//go:build darwin

package lib

import (
	"fmt"
	"net"
	"syscall"
)

const deviceBindSupported = true

// bindToInterface returns a dialer Control hook that pins the socket to the
// interface with IP_BOUND_IF / IPV6_BOUND_IF, the macOS equivalent of Linux's
// SO_BINDTODEVICE.
func bindToInterface(iface *net.Interface) func(network, address string, c syscall.RawConn) error {
	return func(network, _ string, c syscall.RawConn) error {
		var opErr error
		if err := c.Control(func(fd uintptr) {
			switch network {
			case "tcp6", "udp6", "ip6":
				opErr = syscall.SetsockoptInt(int(fd), syscall.IPPROTO_IPV6, syscall.IPV6_BOUND_IF, iface.Index)
			default:
				opErr = syscall.SetsockoptInt(int(fd), syscall.IPPROTO_IP, syscall.IP_BOUND_IF, iface.Index)
			}
		}); err != nil {
			return err
		}
		if opErr != nil {
			return fmt.Errorf("bind socket to %s: %w", iface.Name, opErr)
		}
		return nil
	}
}
