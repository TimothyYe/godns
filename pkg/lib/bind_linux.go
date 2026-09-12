//go:build linux

package lib

import (
	"fmt"
	"net"
	"syscall"
)

const deviceBindSupported = true

// bindToInterface returns a dialer Control hook that binds the socket to the
// device with SO_BINDTODEVICE, so the kernel routes through that interface
// regardless of what the main routing table prefers. Kernels older than 5.7
// require CAP_NET_RAW to set this option.
func bindToInterface(iface *net.Interface) func(network, address string, c syscall.RawConn) error {
	return func(_, _ string, c syscall.RawConn) error {
		var opErr error
		if err := c.Control(func(fd uintptr) {
			opErr = syscall.BindToDevice(int(fd), iface.Name)
		}); err != nil {
			return err
		}
		if opErr != nil {
			return fmt.Errorf("bind socket to %s (kernels before 5.7 need CAP_NET_RAW): %w", iface.Name, opErr)
		}
		return nil
	}
}
