//go:build !linux && !darwin

package lib

import (
	"net"
	"syscall"
)

const deviceBindSupported = false

// bindToInterface is a no-op on platforms without a device-bind socket
// option. The dialer still binds its source address to the interface, which
// is enough when the host's routing already sends that source out the
// intended link.
func bindToInterface(_ *net.Interface) func(network, address string, c syscall.RawConn) error {
	return nil
}
